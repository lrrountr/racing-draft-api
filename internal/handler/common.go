package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/config"
)

func AttachHandler(conf config.Config, r *gin.Engine) {
	//Public health endpoints
	r.GET("/", ok)
	r.GET("/health", ok)

	//Seasons endpoints
	seasons := r.Group("/api/seasons")
	seasons.POST("", CreateNewSeason)
	seasons.GET("", ListSeasons)
	seasons.POST("/:id", UpdateSeason)
	seasons.GET("/:id", GetSeason)
	seasons.DELETE("/:id", DeleteSeason)

	//Users endpoints
	users := r.Group("/api/users")
	users.POST("", CreateUser)
	users.GET("", ListUsers)
	users.POST("/:id", UpdateUser)
	users.GET("/:id", GetUser)
	users.DELETE("/:id", DeleteUser)
}
