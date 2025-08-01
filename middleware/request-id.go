package middleware

import (
	"context"
	"done-hub/common/logger"
	"done-hub/common/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := utils.GetTimeString() + utils.GetRandomString(8)
		c.Set(logger.RequestIdKey, id)
		c.Set("requestStartTime", time.Now())
		ctx := context.WithValue(c.Request.Context(), logger.RequestIdKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(logger.RequestIdKey, id)
		c.Next()
	}
}
