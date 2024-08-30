package model

import (
	"ctx"
)

func (c DBClient) CreateNewSeason(ctx contetx.Context, in CreateNewSeasonRequest) (out CreateNewSeasonResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.RollBack(ctx)

	return out, tx.Commit(ctx)
}
