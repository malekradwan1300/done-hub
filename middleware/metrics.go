package middleware

import (
	"done-hub/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MetricsWithBasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetOrDefault("metrics.user", "")
		password := utils.GetOrDefault("metrics.password", "")

		reqUser, reqPassword, hasAuth := c.Request.BasicAuth()
		if user == "" || password == "" || !hasAuth || reqUser != user || reqPassword != password {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}
