package handler

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	requestUriKey    = "requestUri"
	requestHostKey   = "requestHost"
	requestMethodKey = "requestMethod"
)

func ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
	})
}

func OK(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, body)
}

func InternalServerError(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusInternalServerError, "Internal Server Error - %s", callerMsg, logDetailsAndArgs...)
}

func NotAcceptable(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusNotAcceptable, "Not Acceptable - %s", callerMsg, logDetailsAndArgs...)
}

func Unauthorized(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusUnauthorized, "Unauthorized - %s", callerMsg, logDetailsAndArgs...)
}

func Forbidden(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusUnauthorized, "Forbidden - %s", callerMsg, logDetailsAndArgs...)
}

func NotFound(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusNotFound, "Not Found - %s", callerMsg, logDetailsAndArgs...)
}

func NotProcessable(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusUnprocessableEntity, "Not Processable - %s", callerMsg, logDetailsAndArgs...)
}

func ServiceUnavailable(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusServiceUnavailable, "Service Unavailable - %s", callerMsg, logDetailsAndArgs...)
}

func NotImplemented(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusNotImplemented, "Not Implemented - %s", callerMsg, logDetailsAndArgs...)
}

func StatusConflict(c *gin.Context, callerMsg string, logDetailsAndArgs ...interface{}) {
	respond(c, http.StatusConflict, "Conflict - %s", callerMsg, logDetailsAndArgs...)
}

func respond(c *gin.Context, status int, template string, callerMsg string, logDetailsAndArgs ...interface{}) {
	logMsg := fmt.Sprintf(template, callerMsg)
	var fields *log.Fields
	if len(logDetailsAndArgs) > 0 {
		first := logDetailsAndArgs[0]
		switch m := first.(type) {
		case string:
			logMsg = fmt.Sprintf("%s: %s", logMsg, fmt.Sprintf(m, logDetailsAndArgs[1:]...))
		case error:
			logMsg = fmt.Sprintf("%s: %s", logMsg, m.Error())
		case log.Fields:
			f, ok := logDetailsAndArgs[0].(log.Fields)
			if ok {
				fields = &f
			}
		default:
			log.Errorf("developer error: first argument of logDetailsAndArgs is NOT a string or error")
			goto RESPONSE
		}
	}
	if fields == nil {
		log.Error(logMsg)
	} else {
		(*fields)[requestUriKey] = html.EscapeString(c.Request.RequestURI)
		(*fields)[requestHostKey] = c.Request.Host
		(*fields)[requestMethodKey] = c.Request.Method
		log.WithFields(*fields).Errorf(logMsg)
	}
RESPONSE:
	c.JSON(status, gin.H{
		"msg": fmt.Sprintf(template, callerMsg),
	})
}
