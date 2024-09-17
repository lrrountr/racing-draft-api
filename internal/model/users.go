package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	UUID  string
	Name  string
	Email string
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type CreateUserResponse struct {
	UUID string
}

const (
	createUserSQL = `
	INSERT INTO users
		(uuid, name, email)
	VALUES
		($1, $2, $3)
	ON CONFLICT(email)
	DO UPDATE SET name = $2
`
)

func (c DBClient) CreateUser(ctx context.Context, in CreateUserRequest) (out CreateUserResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, fmt.Errorf("failed to create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	uuid := "user-" + uuid.NewString()

	_, err = tx.Exec(ctx, createUserSQL, uuid, in.Name, in.Email)
	if err != nil {
		return out, fmt.Errorf("failed to create/update user: %w", err)
	}

	out.UUID = uuid
	return out, tx.Commit(ctx)
}

const (
	fetchUserByIdSQL = `
	SELECT 
		uuid, name, email
	FROM
		users
	WHERE uuid = $1
	LIMIT 1
`
)

type FetchUserByIdRequest struct {
	UUID string
}

type FetchUserByIdResponse = User

func (c DBClient) FetchUserById(ctx context.Context, in FetchUserByIdRequest) (out FetchUserByIdResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, fmt.Errorf("failed to create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, fetchUserByIdSQL,
		in.UUID,
	)
	if err != nil {
		return out, fmt.Errorf("query failed: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(
			&out.UUID,
			&out.Name,
			&out.Email,
		)
	}
	defer rows.Close()

	return out, tx.Commit(ctx)
}
