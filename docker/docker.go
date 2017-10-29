package docker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

// BuildImage takes a BuildConfig and generates a docker image.
func BuildImage(ctx context.Context, bc BuildConfig) error {

	if err := bc.Validate(); err != nil {
		return err
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	segs := strings.Split(*bc.Directory, "/")
	tmp := fmt.Sprintf("/tmp/:%s", segs[len(segs)-1])

	tar := archivex.TarFile{}
	tar.Create(tmp)
	tar.AddAll(*bc.Directory, false)
	tar.Close()

	c, err := os.Open(fmt.Sprintf("%s.tar", tmp))
	if err != nil {
		return err
	}
	defer c.Close()

	resp, err := cli.ImageBuild(ctx, c, bc.ToOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListImages returns all the available docker images.
func ListImages(ctx context.Context, all bool) ([]types.ImageSummary, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return cli.ImageList(ctx, types.ImageListOptions{
		All: all,
	})
}
