package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/config"
)

func StartServer(config config.Config) error {
	e := gin.New()

	addr := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	return e.Run(addr)
}
