package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthEndpoints(groups *RouterGroups) {
	groups.public.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up and running",
		})
	})
}
