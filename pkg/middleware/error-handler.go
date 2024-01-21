package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// simple error middleware that returns the occurred errors to the user and logs them
// FIXME: consider not exposing internal errors to the user to mitigate security risk
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		clientIp := c.ClientIP()
		var errors []string
		for _, err := range c.Errors {
			errS := err.Error()
			errors = append(errors, errS)
			log.Printf(`Error: "%s" occurred in url "%s" by IP "%s"`, errS, c.Request.URL, clientIp)
		}

		if len(errors) > 0 {
			// status -1 doesn't overwrite existing status code
			c.JSON(-1, gin.H{"errors": errors})
		}
	}
}
