package corel

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type requestkey string

const RequestIDKey requestkey = "requestId"

func GetRequestIdFromContext(ctx context.Context) string {
	if gc, ok := ctx.(*gin.Context); ok {
		if val, exist := gc.Get(string(RequestIDKey)); exist {
			return val.(string)
		}
	}
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return "no-requestid-in-context"
}

var DefaultGinHandlers = []gin.HandlerFunc{
	func(c *gin.Context) {
		rid := xid.New().String()
		c.Set(string(RequestIDKey), rid)
		c.Header(string(RequestIDKey), rid)
	},
}

func CreateNewContext() context.Context {
	return context.WithValue(context.Background(), RequestIDKey, xid.New().String())
}
