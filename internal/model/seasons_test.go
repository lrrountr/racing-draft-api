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
	_, err := Client.CreateNewSeason(ctx,
		CreateNewSeasonRequest{
			Name:         name,
			RacingSeries: reacingSeries,
			Year:         year,
			Active:       active,
		},
	)
	assert.NilError(t, err)
}
