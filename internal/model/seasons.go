package model

import "context"

type CreateNewSeasonRequest struct {
}

type CreateNewSeasonResponse struct {
}

func (c DBClient) CreateNewSeason(ctx context.Context, in CreateNewSeasonRequest) (out CreateNewSeasonResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx)

	return out, tx.Commit(ctx)
}
