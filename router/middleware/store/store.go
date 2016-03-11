package store

import (
	"github.com/drone/drone/store"
	"github.com/drone/drone/store/datastore"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
)

var (
	driver     = envflag.String("DATABASE_DRIVER", "sqlite3", "")
	datasource = envflag.String("DATABASE_DATASOURCE", "drone.sqlite", "")
)

func Load() gin.HandlerFunc {
	store_ := datastore.New(*driver, *datasource)

	logrus.Infof("using database driver %s", *driver)
	logrus.Infof("using database config %s", *datasource)

	return func(c *gin.Context) {
		store.ToContext(c, store_)
		c.Next()
	}
}
