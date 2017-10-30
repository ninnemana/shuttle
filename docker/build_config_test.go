package docker_test

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/ninnemana/shuttle/docker"
)

func TestBuildConfig_Validate(t *testing.T) {
	dir := "/tmp"
	type fields struct {
		Directory  *string
		Dockerfile *string
		Remove     *bool
		Tags       []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Empty Directory",
			fields{
				Directory: nil,
			},
			true,
		},
		{
			"Success",
			fields{
				Directory: &dir,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := docker.BuildConfig{
				Directory:  tt.fields.Directory,
				Dockerfile: tt.fields.Dockerfile,
				Remove:     tt.fields.Remove,
				Tags:       tt.fields.Tags,
			}
			if err := bc.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("BuildConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildConfig_ToOptions(t *testing.T) {
	emptyDocker := ""
	override := "dev.Dockerfile"
	remove := true
	type fields struct {
		Directory  *string
		Dockerfile *string
		Remove     *bool
		Tags       []string
	}
	tests := []struct {
		name   string
		fields fields
		want   types.ImageBuildOptions
	}{
		{
			"Default Dockerfile",
			fields{},
			types.ImageBuildOptions{
				Tags:           nil,
				SuppressOutput: false,
				Remove:         false,
				ForceRemove:    false,
				PullParent:     true,
				Dockerfile:     docker.DefaultDockerfile,
			},
		},
		{
			"Empty Dockerfile",
			fields{
				Dockerfile: &emptyDocker,
			},
			types.ImageBuildOptions{
				Tags:           nil,
				SuppressOutput: false,
				Remove:         false,
				ForceRemove:    false,
				PullParent:     true,
				Dockerfile:     docker.DefaultDockerfile,
			},
		},
		{
			"Override Dockerfile",
			fields{
				Dockerfile: &override,
			},
			types.ImageBuildOptions{
				Tags:           nil,
				SuppressOutput: false,
				Remove:         false,
				ForceRemove:    false,
				PullParent:     true,
				Dockerfile:     override,
			},
		},
		{
			"With Remove",
			fields{
				Remove: &remove,
			},
			types.ImageBuildOptions{
				Tags:           nil,
				SuppressOutput: false,
				Remove:         remove,
				ForceRemove:    remove,
				PullParent:     true,
				Dockerfile:     docker.DefaultDockerfile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := docker.BuildConfig{
				Directory:  tt.fields.Directory,
				Dockerfile: tt.fields.Dockerfile,
				Remove:     tt.fields.Remove,
				Tags:       tt.fields.Tags,
			}
			if got := bc.ToOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildConfig.ToOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
