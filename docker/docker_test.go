package docker_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/ninnemana/shuttle/docker"
)

const (
	validDocker = `
	FROM alpine:3.6
	RUN apk add --no-cache bash gawk sed grep bc coreutils
	`
)

var (
	env = map[string]string{}
)

func TestMain(m *testing.M) {
	env = envToMap()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working director: %s", err)
		return
	}

	f, err := os.Create(fmt.Sprintf("%s/Dockerfile", dir))
	if err != nil {
		log.Fatalf("Failed to create temporary dockerfile: %s", err)
		return
	}
	_, err = f.WriteString(validDocker)
	if err != nil {
		log.Fatalf("Failed to write to temporary dockerfile: %s", err)
		return
	}

	f.Close()

	m.Run()

	err = os.Remove(fmt.Sprintf("%s/Dockerfile", dir))
	if err != nil {
		log.Fatalf("Failed to remove to temporary dockerfile: %s", err)
		return
	}
}

func TestBuildImage(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	type args struct {
		ctx context.Context
		bc  docker.BuildConfig
	}
	tests := []struct {
		name    string
		args    args
		envs    map[string]string
		wantErr bool
	}{
		{
			"Success",
			args{
				context.Background(),
				docker.BuildConfig{
					Directory: &dir,
					Tags:      []string{"shutle-test:latest"},
				},
			},
			nil,
			false,
		},
		{
			"Bad Docker Client",
			args{
				context.Background(),
				docker.BuildConfig{
					Directory: &dir,
					Tags:      []string{"shutle-test:latest"},
				},
			},
			map[string]string{
				"DOCKER_CERT_PATH": "invalid/path",
			},
			true,
		},
		{
			"NilDirectory",
			args{
				context.Background(),
				docker.BuildConfig{},
			},
			nil,
			true,
		},
	}
	defer mapToEnv(env)
	for _, tt := range tests {
		mapToEnv(env)
		mapToEnv(tt.envs)
		t.Run(tt.name, func(t *testing.T) {
			if err := docker.BuildImage(tt.args.ctx, tt.args.bc); (err != nil) != tt.wantErr {
				t.Errorf("BuildImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestListImages(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		all bool
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		filters map[string]string
// 		wantErr bool
// 	}{
// 		{
// 			"success",
// 			args{
// 				context.Background(),
// 				true,
// 			},
// 			nil,
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		os.Setenv("DOCKER_CERT_PATH", "")
// 		t.Run(tt.name, func(t *testing.T) {
// 			if _, err := docker.ListImages(tt.args.ctx, tt.args.all); (err != nil) != tt.wantErr {
// 				t.Errorf("ListImages() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func envToMap() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitAfterN(e, "=", 2)
		env[kv[0]] = kv[1]
	}

	return env
}

func mapToEnv(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}

func TestListImages(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts *docker.FilterOptions
	}
	all := true
	name := "shutle-test"
	latest := "latest"
	tests := []struct {
		name    string
		args    args
		want    []types.ImageSummary
		wantErr bool
	}{
		{
			"Success",
			args{
				context.Background(),
				&docker.FilterOptions{
					All:  &all,
					Name: &name,
					Tag:  &latest,
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		os.Setenv("DOCKER_CERT_PATH", "")
		t.Run(tt.name, func(t *testing.T) {
			_, err := docker.ListImages(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
