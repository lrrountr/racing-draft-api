package model

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestUserBasics(t *testing.T) {
	ctx, cls := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cls()

	name := uuid.NewString()
	email := uuid.NewString()
	user, err := Client.CreateUser(ctx,
		CreateUserRequest{
			Name:  name,
			Email: email,
		},
	)
	assert.NilError(t, err)

	resp, err := Client.FetchUserById(ctx,
		FetchUserByIdRequest{
			user.UUID,
		},
	)
	assert.NilError(t, err)
	assert.Equal(t, resp.Name, name)
	assert.Equal(t, resp.Email, email)
}
