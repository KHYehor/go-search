package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UniqueRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		c.Set("request-id", id.String())

		c.Next()
	}
}
