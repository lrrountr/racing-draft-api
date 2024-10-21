package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lrrountr/racing-draft-api/internal/clients"
	"github.com/lrrountr/racing-draft-api/internal/model"
)

type CreateNewSeasonRequest struct {
	Name         string `json:"name"`
	RacingSeries string `json:"racing_series"`
	Year         int    `json:"year"`
	Active       bool   `json:"active"`
}

type CreateNewSeasonResponse struct {
}

func CreateNewSeason(c *gin.Context) {
	b := CreateNewSeasonRequest{}
	if BindJSONOrAbort(c, &b) {
		return
	}

	db := clients.LoadDatabase(c)
	_, err := db.CreateNewSeason(c, model.CreateNewSeasonRequest{
		Name:         b.Name,
		RacingSeries: b.RacingSeries,
		Year:         b.Year,
		Active:       b.Active,
	})
	if err != nil {
		InternalServerError(c, "failed to create season", err)
		return
	}

	OK(c, gin.H{
		"msg": "OK",
	})
}

func ListSeasons(c *gin.Context) {
	NotImplemented(c, "come back soon")
}

func UpdateSeason(c *gin.Context) {
	NotImplemented(c, "come back soon")
}

func GetSeason(c *gin.Context) {
	NotImplemented(c, "come back soon")
}

func DeleteSeason(c *gin.Context) {
	NotImplemented(c, "come back soon")
}
