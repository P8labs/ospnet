package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type Client struct {
	cli *client.Client
}

type RunOptions struct {
	Labels map[string]string
	Image  string
	Name   string
	Port   int
	Env    []string
}

type Container struct {
	ID      string
	Name    string
	Image   string
	State   container.ContainerState
	Running bool
}

type ContainerDetails struct {
	ID      string
	Name    string
	Image   string
	State   container.ContainerState
	Running bool
}

func New() (*Client, error) {
	cli, err := client.New(
		client.FromEnv,
		client.WithUserAgent("ospnet-agent"),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client init failed: %w", err)
	}

	return &Client{cli: cli}, nil
}

func (c *Client) PullImage(ctx context.Context, imageRef string) error {
	if strings.TrimSpace(imageRef) == "" {
		return fmt.Errorf("image is required")
	}

	reader, err := c.cli.ImagePull(ctx, imageRef, client.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("image pull failed: %w", err)
	}
	defer reader.Close()

	_, _ = io.Copy(io.Discard, reader)
	return nil
}

func (c *Client) RunContainer(ctx context.Context, opts RunOptions) (string, error) {
	if opts.Image == "" || opts.Name == "" || opts.Port == 0 {
		return "", fmt.Errorf("invalid run options")
	}

	existing, err := c.FindContainersByName(ctx, opts.Name)
	if err == nil && len(existing) > 0 {
		return "", fmt.Errorf("container with name %s already exists", opts.Name)
	}

	if err := c.PullImage(ctx, opts.Image); err != nil {
		return "", err
	}

	containerPort := network.MustParsePort("80/tcp")
	containerCreateOptions := &client.ContainerCreateOptions{
		Image: opts.Image,
		Config: &container.Config{
			Labels: opts.Labels,
			ExposedPorts: network.PortSet{
				containerPort: struct{}{},
			},
			Env: opts.Env,
		},
		HostConfig: &container.HostConfig{
			PortBindings: network.PortMap{
				containerPort: []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
						HostPort: strconv.Itoa(opts.Port),
					},
				},
			},
		},
	}

	resp, err := c.cli.ContainerCreate(ctx, *containerCreateOptions)
	if err != nil {
		return "", fmt.Errorf("create container: %w", err)
	}

	_, err = c.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})

	if err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}

	return resp.ID, nil
}

func (c *Client) StopContainer(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("container id required")
	}

	timeout := 10

	if _, err := c.cli.ContainerStop(ctx, id, client.ContainerStopOptions{
		Timeout: &timeout,
	}); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}

	if _, err := c.cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("remove container: %w", err)
	}

	return nil
}

func (c *Client) ListContainers(ctx context.Context) ([]Container, error) {
	list, err := c.cli.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	var result []Container

	for _, ctr := range list.Items {
		name := ""
		if len(ctr.Names) > 0 {
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}

		result = append(result, Container{
			ID:      ctr.ID,
			Name:    name,
			Image:   ctr.Image,
			State:   ctr.State,
			Running: ctr.State == "running",
		})
	}

	return result, nil
}

func (c *Client) InspectContainer(ctx context.Context, id string) (ContainerDetails, error) {
	if id == "" {
		return ContainerDetails{}, fmt.Errorf("container id required")
	}

	ctr, err := c.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return ContainerDetails{}, fmt.Errorf("inspect container: %w", err)
	}

	return ContainerDetails{
		ID:      ctr.Container.ID,
		Name:    strings.TrimPrefix(ctr.Container.Name, "/"),
		Image:   ctr.Container.Config.Image,
		State:   ctr.Container.State.Status,
		Running: ctr.Container.State.Running,
	}, nil
}

func (c *Client) FindContainersByName(ctx context.Context, name string) ([]Container, error) {
	args := client.Filters{}
	args.Add("name", name)

	list, err := c.cli.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return nil, fmt.Errorf("find containers: %w", err)
	}

	var result []Container

	for _, ctr := range list.Items {
		n := ""
		if len(ctr.Names) > 0 {
			n = strings.TrimPrefix(ctr.Names[0], "/")
		}

		result = append(result, Container{
			ID:      ctr.ID,
			Name:    n,
			Image:   ctr.Image,
			State:   ctr.State,
			Running: ctr.State == "running",
		})
	}

	return result, nil
}

func ParsePortBinding(raw string) (int, error) {
	_, portValue, err := net.SplitHostPort(raw)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(portValue)
}
