package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/clients"
	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/model"
)

func AttachConfig(c *gin.Context, conf config.Config) {
	clients.AttachConfig(c, conf)
}

func AttachDatabase(c *gin.Context, db model.DBClient) {
	clients.AttachDatabase(c, db)
}
