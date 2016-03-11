package docker

import (
	"io"

	"github.com/samalba/dockerclient"
)

var (
	LogOpts = &dockerclient.LogOptions{
		Stdout: true,
		Stderr: true,
	}

	LogOptsTail = &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	}
)

func NewClient() (dockerclient.Client, error) {
	return dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
}

func MustClient() dockerclient.Client {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	return client
}

// Run creates the docker container, pulling images if necessary, starts
// the container and blocks until the container exits, returning the exit
// information.
func Run(client dockerclient.Client, conf *dockerclient.ContainerConfig, name string) (*dockerclient.ContainerInfo, error) {
	info, err := RunDaemon(client, conf, name)
	if err != nil {
		return nil, err
	}

	<-client.Wait(info.Id)

	return client.InspectContainer(info.Id)
}

// RunDaemon creates the docker container, pulling images if necessary, starts
// the container and returns the container information. It does not wait for
// the container to exit.
func RunDaemon(client dockerclient.Client, conf *dockerclient.ContainerConfig, name string) (*dockerclient.ContainerInfo, error) {

	// attempts to create the contianer
	id, err := client.CreateContainer(conf, name, nil)
	if err != nil {
		// and pull the image and re-create if that fails
		err = client.PullImage(conf.Image, nil)
		if err != nil {
			return nil, err
		}
		id, err = client.CreateContainer(conf, name, nil)
		if err != nil {
			client.RemoveContainer(id, true, true)
			return nil, err
		}
	}

	// fetches the container information
	info, err := client.InspectContainer(id)
	if err != nil {
		client.RemoveContainer(id, true, true)
		return nil, err
	}

	// starts the container
	err = client.StartContainer(id, &conf.HostConfig)
	if err != nil {
		client.RemoveContainer(id, true, true)
		return nil, err
	}

	return info, err
}

// Tail tails the container logs to the io.Writer for the
// named Docker container.
func Tail(client dockerclient.Client, w io.Writer, name string) error {
	rc, err := client.ContainerLogs(name, LogOptsTail)
	if err != nil {
		return err
	}
	StdCopy(w, w, rc)
	rc.Close()
	return nil
}

// Delete stops and deletes the container, including any
// container volumes.
func Delete(client dockerclient.Client, name string) {
	client.StopContainer(name, 30)
	client.KillContainer(name, "9")
	client.RemoveContainer(name, true, true)
}
