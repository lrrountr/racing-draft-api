package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateNewSeasonRequest struct {
	Name         string
	RacingSeries string
	Year         int
	Active       bool
}

type CreateNewSeasonResponse struct {
	UUID string
}

const (
	createNewSeasonSQL = `
	INSERT INTO seasons
		(uuid, name, racing_series, year, active)
	VALUES
		($1, $2, $3, $4, $5)
	ON CONFLICT (uuid)
	DO UPDATE SET name = $2, racing_series = $3, year= $4, active = $5
`
)

func (c DBClient) CreateNewSeason(ctx context.Context, in CreateNewSeasonRequest) (out CreateNewSeasonResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, fmt.Errorf("failed to create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	uuid := "season-" + uuid.NewString()
	_, err = tx.Exec(ctx, createNewSeasonSQL,
		uuid,
		in.Name,
		in.RacingSeries,
		in.Year,
		in.Active,
	)
	if err != nil {
		return out, fmt.Errorf("failed to create/update season: %w", err)
	}

	out.UUID = uuid
	return out, tx.Commit(ctx)
}

type Season struct {
	UUID         string
	Name         string
	RacingSeries string
	Year         int
	Active       bool
}

type ListSeasonsRequest struct {
	RacingSeries string
	Limit        int
	Offset       int
}

type ListSeasonsResponse struct {
	PageInfo PageInfo
	Seasons  []Season
}

const (
	listSeasonsSQL = `
	SELECT
		uuid,
		name,
		racing_series,
		year,
		active
	FROM
		seasons
	WHERE
		racing_series = $1 OR $1 IS NULL
	LIMIT $2
	OFFSET $3	
`
	listSeasonsTotalSQL = `
	SELECT
		COUNT(uuid)
	FROM
		seasons
	WHERE
		racing_series = $1 OR $1 IS NULL
	LIMIT 1
`
)

func (c DBClient) ListSeasons(ctx context.Context, in ListSeasonsRequest) (out ListSeasonsResponse, err error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return out, fmt.Errorf("failed to create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, listSeasonsSQL,
		in.RacingSeries,
		in.Limit,
		in.Offset,
	)
	if err != nil {
		return out, fmt.Errorf("query failed: %w", err)
	}
	for rows.Next() {
		s := Season{}
		err = rows.Scan(
			&s.UUID,
			&s.Name,
			&s.RacingSeries,
			&s.Year,
			&s.Active,
		)
		out.Seasons = append(out.Seasons, s)
	}
	defer rows.Close()

	row := tx.QueryRow(ctx, listSeasonsTotalSQL)
	total := 0
	err = row.Scan(&total)
	if err != nil {
		return out, fmt.Errorf("failed to scan total rows: %w", err)
	}
	out.PageInfo.Update(total, len(out.Seasons), in.Offset)

	return out, tx.Commit(ctx)
}
