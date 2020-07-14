package tests_common

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Container struct {
	Address      string
	Image        string
	Name         string
	Environments []string
	Ports        map[string]string
	Mounts       map[string]string
	Cmd          []string
	id           string
	Inspection   *types.ContainerJSON
}

func (c *Container) Run() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	hasImage := false
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
imagesLoop:
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, fmt.Sprintf("%s:", c.Image)) {
				hasImage = true
				break imagesLoop
			}
		}

	}
	if !hasImage {
		reader, err := cli.ImagePull(ctx, c.Address, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		io.Copy(os.Stderr, reader)
	}

	cli.ContainerRemove(ctx, c.Name, types.ContainerRemoveOptions{
		Force: true,
	})

	pb := nat.PortMap{}
	if c.Ports != nil {
		for k, v := range c.Ports {
			pb[nat.Port(v)] = []nat.PortBinding{
				{
					HostIP:   strings.Split(k, ":")[0],
					HostPort: strings.Split(k, ":")[1],
				},
			}
		}
	}

	mounts := []mount.Mount{}
	for host, target := range c.Mounts {
		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   host,
			Target:   target,
			ReadOnly: false,
		})
	}

	cntr, err := cli.ContainerCreate(ctx, &container.Config{
		Env:   c.Environments,
		Image: c.Image,
		Cmd:   c.Cmd,
	}, &container.HostConfig{
		PortBindings: pb,
		Mounts:       mounts,
	}, nil, c.Name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, cntr.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	c.id = cntr.ID
	s, err := cli.ContainerInspect(ctx, cntr.ID)
	if err != nil {
		panic(err)
	}
	c.Inspection = &s
}

func (c *Container) Close() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	err = cli.ContainerRemove(ctx, c.Name, types.ContainerRemoveOptions{
		Force: true,
	})
}
