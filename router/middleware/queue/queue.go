package queue

import (
	"github.com/drone/drone/queue"
	"github.com/gin-gonic/gin"
)

func Load() gin.HandlerFunc {
	queue_ := queue.New()
	return func(c *gin.Context) {
		queue.ToContext(c, queue_)
		c.Next()
	}
}
