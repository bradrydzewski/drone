package engine

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/drone/drone/bus"
	"github.com/drone/drone/engine/docker"
	"github.com/drone/drone/model"
	"github.com/drone/drone/queue"
	"github.com/drone/drone/store"
	"github.com/drone/drone/stream"
	"github.com/samalba/dockerclient"
	"golang.org/x/net/context"
)

// Poll polls the build queue for build jobs.
func Poll(c context.Context) {
	for {
		pollRecover(c)
	}
}

func pollRecover(c context.Context) {
	defer recover()
	poll(c)
}

func poll(c context.Context) {
	w := queue.Pull(c)

	logrus.Infof("Starting build %s/%s#%d.%d",
		w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number)

	name := cname(w)
	client := docker.FromContext(c)
	rc, wc, err := stream.Create(c, stream.ToKey(w.Job.ID))
	if err != nil {
		logrus.Errorf("Error opening build stream %s/%s#%d.%d. %s",
			w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number, err)
	}

	defer func() {
		wc.Close()
		rc.Close()
		stream.Remove(c, stream.ToKey(w.Job.ID))
		docker.Delete(client, name)
	}()

	// run the build
	stdin, _ := json.Marshal(w)
	args := []string{ /*"--pull", "--cache",*/ "--clone", "--build", "--deploy"}
	args = append(args, "--")
	args = append(args, string(stdin))
	conf := &dockerclient.ContainerConfig{
		Image:      "drone/drone-exec:latest",
		Entrypoint: []string{"/bin/drone-exec"},
		Cmd:        args,
		HostConfig: dockerclient.HostConfig{
			Binds: []string{"/var/run/docker.sock:/var/run/docker.sock"},
		},
		Volumes: map[string]struct{}{
			"/var/run/docker.sock": struct{}{},
		},
	}
	if _, err := docker.RunDaemon(client, conf, name); err != nil {
		logrus.Errorf("Error starting build %s/%s#%d.%d. %s",
			w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number, err)
	}

	go func() {
		err := docker.Tail(client, stream.NewWriter(wc), name)
		if err != nil {
			logrus.Errorf("Error writing build stream %s/%s#%d.%d. %s",
				w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number, err)
		}
		wc.Close()
	}()

	w.Job.Status = model.StatusRunning
	w.Job.Started = time.Now().Unix()

	quitc := make(chan bool, 1)
	eventc := make(chan *bus.Event, 1)
	bus.Subscribe(c, eventc)

	defer func() {
		bus.Unsubscribe(c, eventc)
		quitc <- true
	}()
	go func() {
		for {
			select {
			case event := <-eventc:
				if event.Type == bus.Cancelled && event.Job.ID == w.Job.ID {
					logrus.Infof("Cancel build %s/%s#%d.%d",
						w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number)
					go client.StopContainer(name, 30)
				}
			case <-quitc:
				return
			}
		}
	}()

	// TODO store the started build in the database
	// TODO publish the started build
	store.UpdateJob(c, w.Job)
	//store.Write(c, w.Job, rc)
	bus.Publish(c, bus.NewEvent(bus.Started, w.Repo, w.Build, w.Job))

	// catch the build result
	result := <-client.Wait(name)
	w.Job.ExitCode = result.ExitCode
	w.Job.Finished = time.Now().Unix()

	switch w.Job.ExitCode {
	case 128, 130:
		w.Job.Status = model.StatusKilled
	case 0:
		w.Job.Status = model.StatusSuccess
	default:
		w.Job.Status = model.StatusFailure
	}

	// store the finished build in the database
	logs, _, err := stream.Open(c, stream.ToKey(w.Job.ID))
	if err != nil {
		logrus.Errorf("Error reading build stream %s/%s#%d.%d",
			w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number)
	}
	defer func() {
		if logs != nil {
			logs.Close()
		}
	}()
	if err := store.WriteLog(c, w.Job, logs); err != nil {
		logrus.Errorf("Error persisting build stream %s/%s#%d.%d",
			w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number)
	}
	if logs != nil {
		logs.Close()
	}

	// TODO publish the finished build
	store.UpdateJob(c, w.Job)
	bus.Publish(c, bus.NewEvent(bus.Finished, w.Repo, w.Build, w.Job))

	logrus.Infof("Finished build %s/%s#%d.%d",
		w.Repo.Owner, w.Repo.Name, w.Build.Number, w.Job.Number)
}

// cname generates a semi-unique deterministic container
// name for a build.
func cname(p *queue.Work) string {
	hash := sha256.New()
	fmt.Fprintln(hash,
		p.Repo.Owner,
		p.Repo.Name,
		p.Build.Number,
		p.Job.Number,
		p.System.Link,
	)
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return fmt.Sprintf("build-%s", sum[:12])
}
