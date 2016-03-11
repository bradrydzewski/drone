package engine

import (
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/drone/drone/engine"
	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
	"github.com/samalba/dockerclient"
)

var once sync.Once

var (
	hosts  = envflag.String("DOCKER_HOST", "unix:///var/run/docker.sock", "")
	certs  = envflag.String("DOCKER_CERT_PATH", "", "")
	verify = envflag.Bool("DOCKER_TLS_VERIFY", false, "")
	procs  = envflag.Int("DOCKER_MAX_PROCS", 2, "")
	agent  = envflag.String("AGENT_TOKEN", "", "")
)

func Load(c *gin.Context) {
	once.Do(func() {
		prepareDocker(c)
	})

	c.Next()
}

// prepareDocker is a simple helper function that parses the docker
// daemon configuration and launches the docker engines to poll
// and begin processing builds.
func prepareDocker(c *gin.Context) {
	if *agent != "" {
		logrus.Infof("Running in agent mode")
		return
	}
	for _, host := range strings.Split(*hosts, " ") {
		for i := 0; i < *procs; i++ {

			client, err := dockerclient.NewDockerClient(host, nil)
			if err != nil {
				logrus.Fatalf("Unable to initialize Docker deamon %s. %s", host, err)
			}
			_, err = client.Info()
			if err != nil {
				logrus.Fatalf("Unable to connect with Docker deamon %s. %s", host, err)
			}
			logrus.Infof("Registered Docker daemon %s. ", host)

			c := c.Copy()
			c.Set("docker", client)
			go engine.Poll(c)
		}
	}
}
