package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func BindQueryOrAbort(c *gin.Context, in interface{}) (doAbort bool) {
	err := c.ShouldBindQuery(in)
	if err != nil {
		log.Errorf("failed to bind query parameters: %s", err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"msg": fmt.Sprintf("Not Acceptable - failed to bind query parameters: %s", err),
		})
		return true
	}
	return false
}

func BindURIOrAbort(c *gin.Context, in interface{}) (doAbort bool) {
	err := c.ShouldBindUri(in)
	if err != nil {
		log.Errorf("failed to bind URI parameters: %s", err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"msg": fmt.Sprintf("Not Acceptable - failed to bind URI parameters: %s", err),
		})
		return true
	}
	return false
}

func BindJSONOrAbort(c *gin.Context, in interface{}) (doAbort bool) {
	err := c.ShouldBindJSON(in)
	if err != nil {
		log.Errorf("failed to bind body parameters: %s", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"msg": "Not Acceptable - failed to bind body parameters",
		})
		return true
	}
	return false
}
