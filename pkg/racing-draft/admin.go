package racing_draft

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lrrountr/racing-draft-api/internal/handler"
)

type (
	CreateNewSeasonRequest  = handler.CreateNewSeasonRequest
	CreateNewSeasonResponse = handler.CreateNewSeasonResponse
)

func (c *client) CreateNewSeason(ctx context.Context, in CreateNewSeasonRequest) (out CreateNewSeasonResponse, err error) {
	path := "/api/seasons"

	body, err := c.createBody(in)
	if err != nil {
		return out, fmt.Errorf("illegal JSOn chars in request (?): %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := c.doRequest(ctx, http.MethodPost, path, body, headers, nil)
	if err != nil {
		return out, fmt.Errorf("could not execute request:  %w", err)
	}
	defer resp.Body.Close()

	err = c.parseResp(resp, &out)
	if err != nil {
		return out, fmt.Errorf("could no parse response: %w", err)
	}

	return out, nil
}
