package stream

import (
	"github.com/drone/drone/stream"
	"github.com/gin-gonic/gin"
)

func Load() gin.HandlerFunc {
	stream_ := stream.New()
	return func(c *gin.Context) {
		stream.ToContext(c, stream_)
		c.Next()
	}
}
