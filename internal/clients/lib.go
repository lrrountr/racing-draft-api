package clients

import (
	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/model"
)

const (
	ConfigContextKey   = "ConfigContext"
	DatabaseContextKey = "DatabaseContext"
)

func AttachConfig(c *gin.Context, conf config.Config) {
	c.Set(ConfigContextKey, conf)
}

func LoadConfig(c *gin.Context) config.Config {
	conf, ok := c.Get(ConfigContextKey)
	if !ok {
		panic("Config no set. This is a dev error.")
	}
	return conf.(config.Config)
}

func AttachDatabase(c *gin.Context, db model.DBClient) {
	c.Set(DatabaseContextKey, db)
}

func LoadDatabase(c *gin.Context) model.DBClient {
	db, ok := c.Get(DatabaseContextKey)
	if !ok {
		panic("Da")
	}
	return db.(model.DBClient)
}
