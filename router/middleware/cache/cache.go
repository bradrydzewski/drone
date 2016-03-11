package cache

import (
	"time"

	"github.com/drone/drone/cache"

	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
)

var ttl = envflag.Duration("CACHE_TTL", time.Minute*15, "")

func Load() gin.HandlerFunc {
	cache_ := cache.NewTTL(*ttl)
	return func(c *gin.Context) {
		cache.ToContext(c, cache_)
		c.Next()
	}
}
