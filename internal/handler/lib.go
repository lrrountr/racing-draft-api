package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/model"
)

var AllowHeaders = []string{
	"Origin",
	"Content-Length",
	"Access-Control-Allow-Origin",
	"Content-Type",
	"Authorization",
}

func StartServer(config config.Config) error {
	e := gin.New()

	skipLogs := []string{
		"/",
	}
	handlers := []gin.HandlerFunc{
		gin.LoggerWithWriter(gin.DefaultWriter, skipLogs...),
		gin.Recovery(),
	}
	e.Use(handlers...)

	corsConf := cors.DefaultConfig()
	corsConf.AllowOriginFunc = func(origin string) bool { return true }
	corsConf.AllowCredentials = true
	corsConf.AllowHeaders = AllowHeaders
	e.Use(cors.New(corsConf))

	dbClient, err := model.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create DB client: %w", err)
	}

	e.Use(func(c *gin.Context) {
		AttachConfig(c, config)
		AttachDatabase(c, dbClient)
		c.Next()
	})

	AttachHandler(config, e)
	addr := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	return e.Run(addr)
}

func ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
	})
}

func AttachHandler(conf config.Config, r *gin.Engine) {
	//Public health endpoints
	r.GET("/", ok)
	r.GET("/health", ok)
}
