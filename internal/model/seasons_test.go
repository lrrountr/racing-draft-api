package model

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestSeasonsBasics(t *testing.T) {
	ctx, cls := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cls()

	name := uuid.NewString()
	reacingSeries := uuid.NewString()
	year := rand.Intn(2024)
	active := randBool()
	season, err := Client.CreateNewSeason(ctx,
		CreateNewSeasonRequest{
			Name:         name,
			RacingSeries: reacingSeries,
			Year:         year,
			Active:       active,
		},
	)
	assert.NilError(t, err)

	seasons, err := Client.ListSeasons(ctx,
		ListSeasonsRequest{
			RacingSeries: reacingSeries,
			Limit:        2,
		},
	)
	assert.NilError(t, err)
	assert.Equal(t, len(seasons.Seasons), 1)
	assert.Equal(t, seasons.Seasons[0].UUID, season.UUID)
	assert.Equal(t, seasons.PageInfo.Next, 1)
	assert.Equal(t, seasons.PageInfo.Total, 1)
}
