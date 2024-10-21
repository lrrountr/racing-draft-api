package racing_draft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type RacingDraftClient interface {
	UserClient
	AdminClient
}

type UserClient interface {
}

type AdminClient interface {
	CreateNewSeason(ctx context.Context, req CreateNewSeasonRequest) (out CreateNewSeasonResponse, err error)
}

type ClientOpt func(*client) *client

type client struct {
	http.Client

	baseURL string

	lock *sync.Mutex

	accessToken string
	expiresAt   time.Time
}

func NewClient(opts ...ClientOpt) (RacingDraftClient, error) {
	c := &client{
		lock: new(sync.Mutex),
	}
	for _, opt := range opts {
		c = opt(c)
	}
	return c, nil
}

func WithBaseURL(url string) ClientOpt {
	return func(c *client) *client {
		c.baseURL = url
		return c
	}
}

func WithAccessToken(token string) ClientOpt {
	return func(c *client) *client {
		c.accessToken = token
		c.expiresAt = time.Now().Add(10000 * time.Minute)
		return c
	}
}

func createClient(urlBase, token string) (RacingDraftClient, error) {
	_, err := url.Parse(urlBase)
	if err != nil {
		return nil, fmt.Errorf("Destination not a valid URL base: %w", err)
	}
	return &client{
		baseURL: urlBase,

		lock:        new(sync.Mutex),
		accessToken: token,
		expiresAt:   time.Now().Add(10000 * time.Minute),
	}, nil
}

func (c *client) doRequest(ctx context.Context, method, uri string, body io.Reader, extraHeaders map[string]string, qparams map[string][]string) (resp *http.Response, err error) {
	u, e := url.Parse(c.baseURL)
	if e != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	given, e := url.Parse(uri)
	if e != nil {
		return nil, errorf("invalid request URL", http.StatusInternalServerError, err)
	}
	u.Path = given.Path
	values := given.Query()

	for k, vs := range qparams {
		for _, v := range vs {
			values.Add(k, v)
		}
	}

	u.RawQuery = values.Encode()

	url := u.String()
	req, e := http.NewRequestWithContext(ctx, method, url, body)
	if e != nil {
		return nil, err
	}
	auth, ok := extraHeaders["Authorization"]
	switch {
	case ok && auth == "":
		// The user set Authorization,
		// but set it to ""
		// This means we should not use auth for this header at all.
		delete(extraHeaders, "Authorization")
		break
	default:
		accessToken, err := c.getAccessToken()
		if err != nil {
			return nil, errorf("failed to get access token", http.StatusUnauthorized, err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	for k, v := range extraHeaders {
		req.Header.Add(k, v)
	}

	// If we get passed a body, and the extra headers didn't specify a content type,
	// default to application/json
	if body != nil {
		v := req.Header.Get("Content-Type")
		if v == "" {
			req.Header.Add("Content-Type", "application/json")
		}
	}

	resp, e = c.Do(req)
	if e != nil {
		return nil, errorf("could not execute HTTP request", http.StatusInternalServerError, e)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return resp, c.generateError(resp)
	}
	return resp, nil

}

func (c *client) generateError(resp *http.Response) (err error) {
	out, e := io.ReadAll(resp.Body)
	if e != nil {
		return errorf("could not ready body after code", resp.StatusCode, e)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return errorf("resource not found", resp.StatusCode, respToError(out))
	default:
		return errorf("invalid status code returned", resp.StatusCode, respToError(out))
	}

}

func (c *client) getAccessToken() (string, error) {
	return c.accessToken, nil
}

func errorf(base string, status int, cause error) error {
	return fmt.Errorf("%d - %s: %w", status, base, cause)
}

func (c *client) createBody(in interface{}) (r io.Reader, err error) {
	b := bytes.NewBuffer(make([]byte, 0))
	enc := json.NewEncoder(b)
	e := enc.Encode(in)
	if e != nil {
		return nil, errorf("failed to encode body", http.StatusInternalServerError, e)
	}
	return b, nil
}

func (c *client) parseResp(resp *http.Response, in interface{}) (err error) {
	dec := json.NewDecoder(resp.Body)
	e := dec.Decode(in)
	if e != nil {
		return errorf("failed to decode body", http.StatusInternalServerError, e)
	}
	return nil
}
