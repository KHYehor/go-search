package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()
		v := c.GetString("request-id")
		fmt.Printf("Request %s processed in %v with status %d\n", v, time.Since(t), c.Writer.Status())
	}
}
