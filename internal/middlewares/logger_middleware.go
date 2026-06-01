package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf(
			"[REQUEST] %s %s",
			c.Request.Method,
			c.Request.URL.Path,
		)

		c.Next()

		log.Printf(
			"[RESPONSE] %s %s %d (%v)",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)

		if len(c.Errors) > 0 {
			log.Printf(
				"[ERROR] %s",
				c.Errors.String(),
			)
		}
	}
}
