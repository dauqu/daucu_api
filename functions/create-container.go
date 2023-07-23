package functions

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreateDockerContainer(name string, image string, public_port string, exposed_port string, host_path string, env []string) (container_id string, err error) {

	//Create docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	//Create context
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	//CHeck if docker is running
	_, err = cli.Ping(ctx)
	if err != nil {
		return "", err
	}

	//Create Volume
	// volume, err := cli.VolumeCreate(ctx, volume.VolumeCreateBody{
	// 	Name: name,
	// })
	// if err != nil {
	// 	return "", err
	// }

	//HostConfig
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(exposed_port): []nat.PortBinding{
				{
					HostIP:   "localhost",
					HostPort: public_port,
				},
			},
		},
		//Restart Policy
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		//LogConfig
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{"max-size": "100m"},
		},
		//Mounts
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:        mount.TypeVolume,
		// 		Source:      volume.Name,
		// 		Target:      host_path,
		// 		//Cunsistency backup data
		// 		Consistency: mount.ConsistencyCached,
		// 	},
		// },
	}

	//Config
	config := &container.Config{
		Image: image,
		Env:   env,
		Labels: map[string]string{
			"app": name,
		},
		Hostname: name,
		ExposedPorts: nat.PortSet{
			nat.Port(exposed_port): struct{}{},
		},
	}

	//NetworkConfig
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	gatewayConfig := &network.EndpointSettings{
		Gateway: "gatewayname",
	}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	//Stop container if it is running
	// err = cli.ContainerStop(ctx, name, nil)
	// if err != nil {
	// 	//Check if container is not found
	// 	if err.Error() != "No such container: "+name {
	// 		//Skip error
	// 	} else {
	// 		return "", err
	// 	}
	// }

	// Check if container already exists and remove it
	_, err = cli.ContainerInspect(ctx, name)
	if err == nil {
		// Remove container
		err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{})
		if err != nil {
			return "", err
		}
	}

	// Create container
	newContainer, err := cli.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, name)
	if err != nil {
		return "", err
	}

	// Start container and return container ID
	err = cli.ContainerStart(ctx, newContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	//Return container ID
	return newContainer.ID, nil
}
