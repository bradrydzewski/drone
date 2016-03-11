package controller

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/drone/drone/bus"
	"github.com/drone/drone/model"
	"github.com/drone/drone/router/middleware/session"
	"github.com/drone/drone/store"
	"github.com/drone/drone/stream"

	log "github.com/Sirupsen/logrus"

	"github.com/manucorporat/sse"
)

// GetRepoEvents will upgrade the connection to a Websocket and will stream
// event updates to the browser.
func GetRepoEvents(c *gin.Context) {
	repo := session.Repo(c)
	c.Writer.Header().Set("Content-Type", "text/event-stream")

	eventc := make(chan *bus.Event, 1)
	bus.Subscribe(c, eventc)
	defer func() {
		bus.Unsubscribe(c, eventc)
		close(eventc)
		log.Infof("closed event stream")
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case event := <-eventc:
			if event == nil {
				log.Infof("nil event received")
				return false
			}

			if event.Repo.FullName == repo.FullName {

				var payload = struct {
					model.Build
					Jobs []*model.Job `json:"jobs"`
				}{}
				payload.Build = event.Build
				payload.Jobs, _ = store.GetJobList(c, &event.Build)
				data, _ := json.Marshal(&payload)

				sse.Encode(w, sse.Event{
					Event: "message",
					Data:  string(data),
				})
			}
		case <-c.Writer.CloseNotify():
			return false
		}
		return true
	})
}

func GetStream(c *gin.Context) {

	repo := session.Repo(c)
	buildn, _ := strconv.Atoi(c.Param("build"))
	jobn, _ := strconv.Atoi(c.Param("number"))

	c.Writer.Header().Set("Content-Type", "text/event-stream")

	build, err := store.GetBuildNumber(c, repo, buildn)
	if err != nil {
		log.Debugln("stream cannot get build number.", err)
		c.AbortWithError(404, err)
		return
	}
	job, err := store.GetJobNumber(c, build, jobn)
	if err != nil {
		log.Debugln("stream cannot get job number.", err)
		c.AbortWithError(404, err)
		return
	}

	rc, wc, err := stream.Open(c, stream.ToKey(job.ID))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	defer func() {
		if wc != nil {
			wc.Close()
		}
		if rc != nil {
			rc.Close()
		}
	}()

	go func() {
		<-c.Writer.CloseNotify()
		rc.Close()
	}()

	var line int
	var scanner = bufio.NewScanner(rc)
	for scanner.Scan() {
		line++
		var err = sse.Encode(c.Writer, sse.Event{
			Id:    strconv.Itoa(line),
			Event: "message",
			Data:  scanner.Text(),
		})
		if err != nil {
			break
		}
		c.Writer.Flush()
	}
}

type StreamWriter struct {
	writer gin.ResponseWriter
	count  int
}

func (w *StreamWriter) Write(data []byte) (int, error) {
	var err = sse.Encode(w.writer, sse.Event{
		Id:    strconv.Itoa(w.count),
		Event: "message",
		Data:  string(data),
	})
	w.writer.Flush()
	w.count += len(data)
	return len(data), err
}
