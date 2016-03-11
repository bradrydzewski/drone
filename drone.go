// @APIVersion 1.0.0
// @APITitle Drone API
package main

import (
	"net/http"

	"github.com/drone/drone/router"
	"github.com/drone/drone/router/middleware/bus"
	"github.com/drone/drone/router/middleware/cache"
	"github.com/drone/drone/router/middleware/engine"
	"github.com/drone/drone/router/middleware/header"
	"github.com/drone/drone/router/middleware/queue"
	"github.com/drone/drone/router/middleware/remote"
	"github.com/drone/drone/router/middleware/store"
	"github.com/drone/drone/router/middleware/stream"

	"github.com/Sirupsen/logrus"
	"github.com/ianschenck/envflag"
	_ "github.com/joho/godotenv/autoload"
)

// build revision number populated by the continuous
// integration server at compile time.
var build string

var (
	addr = envflag.String("SERVER_ADDR", ":8000", "")
	cert = envflag.String("SERVER_CERT", "", "")
	key  = envflag.String("SERVER_KEY", "", "")

	debug = envflag.Bool("DEBUG", false, "")
)

func main() {
	envflag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// setup the server and start the listener
	handler := router.Load(
		header.Version,
		bus.Load(),
		stream.Load(),
		cache.Load(),
		queue.Load(),
		store.Load(),
		remote.Load(),
		engine.Load,
	)

	if *cert != "" {
		logrus.Fatal(
			http.ListenAndServeTLS(*addr, *cert, *key, handler),
		)
	} else {
		logrus.Fatal(
			http.ListenAndServe(*addr, handler),
		)
	}
}
