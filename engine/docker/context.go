package docker

import (
	"github.com/samalba/dockerclient"
	"golang.org/x/net/context"
)

const key = "docker"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Bus associated with this context.
func FromContext(c context.Context) dockerclient.Client {
	return c.Value(key).(dockerclient.Client)
}

// ToContext adds the Docker client to this context if it
// supports the Setter interface.
func ToContext(c Setter, d dockerclient.Client) {
	c.Set(key, d)
}
