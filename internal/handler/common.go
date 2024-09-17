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
	seasons.PUT("/:id", UpdateSeason)
	seasons.GET("/:id", GetSeason)
	seasons.DELETE("/:id", DeleteSeason)

	//Users endpoints
	users := r.Group("/api/users")
	users.POST("", CreateUser)
	seasons.GET("", ListUsers)
	seasons.POST("/:id", UpdateUser)
	seasons.GET("/:id", GetUser)
	seasons.DELETE("/:id", DeleteUser)
}
