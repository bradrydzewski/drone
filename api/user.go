package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/drone/drone/cache"
	"github.com/drone/drone/router/middleware/session"
	"github.com/drone/drone/shared/token"
	"github.com/drone/drone/store"
)

func GetSelf(c *gin.Context) {
	c.JSON(200, session.User(c))
}

func GetFeed(c *gin.Context) {
	repos, err := cache.GetRepos(c, session.User(c))
	if err != nil {
		c.String(500, "Error fetching repository list. %s", err)
		return
	}

	feed, err := store.GetUserFeed(c, repos)
	if err != nil {
		c.String(500, "Error fetching feed. %s", err)
	} else {
		c.JSON(200, feed)
	}
}

func GetRepos(c *gin.Context) {
	repos, err := cache.GetRepos(c, session.User(c))
	if err != nil {
		c.String(500, "Error fetching repository list. %s", err)
		return
	}

	repos_, err := store.GetRepoListOf(c, repos)
	if err != nil {
		c.String(500, "Error fetching repository list. %s", err)
	} else {
		c.JSON(http.StatusOK, repos_)
	}
}

func GetRemoteRepos(c *gin.Context) {
	repos, err := cache.GetRepos(c, session.User(c))
	if err != nil {
		c.String(500, "Error fetching repository list. %s", err)
	} else {
		c.JSON(http.StatusOK, repos)
	}
}

func PostToken(c *gin.Context) {
	user := session.User(c)

	token := token.New(token.UserToken, user.Login)
	tokenstr, err := token.Sign(user.Hash)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.String(http.StatusOK, tokenstr)
	}
}
