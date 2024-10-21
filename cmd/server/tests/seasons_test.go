package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	shared "github.com/lrrountr/racing-draft-api/pkg/racing-draft"
)

func TestCreateNewSeason(t *testing.T) {
	ctx, cls := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cls()

	client := getAdminClient(t)

	name := uuid.NewString()
	racingSeries := uuid.NewString()
	year := rand.Intn(2030)
	active := rand.Intn(2) == 1
	_, err := client.CreateNewSeason(ctx, shared.CreateNewSeasonRequest{
		Name:         name,
		RacingSeries: racingSeries,
		Year:         year,
		Active:       active,
	})
	assert.NilError(t, err)
}
