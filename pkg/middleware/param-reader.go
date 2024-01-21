package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: consider further abstracting these functions by making it more generic and having it constructed
// with which value should be read from the param

func IdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		if idString == "" {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("faulty redirect; no id param found"))
			return
		}

		c.Set("id", idString)
		c.Next()
	}
}

func LabelParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		labelString := c.Param("label")
		if labelString == "" {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("faulty redirect; no label param found"))
			return
		}

		c.Set("label", labelString)
		c.Next()
	}
}
