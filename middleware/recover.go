package middleware

import (
	"apiingolang/activity/business/utils/logging"
	"context"

	"github.com/gin-gonic/gin"
)

// Recovery returns a gin.HandlerFunc having recovery solution
func Recovery(l logging.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Recover(c, err)
			}
		}()
		c.Next()
	}
}

func Recover(ctx context.Context, err any) {
	logging.Logger.WriteLogs(ctx, "Panic", logging.ErrorLevel, logging.Fields{"error": err})
}
