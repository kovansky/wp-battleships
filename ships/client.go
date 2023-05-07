package ships

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"net/url"
	"time"
)

const ApiTokenHeader = "X-Auth-Token"

var _ battleships.Client = (*Client)(nil)

type Client struct {
	baseUrl string
	client  http.Client
	log     *zerolog.Logger
	ctx     context.Context
}

func NewClient(ctx context.Context, baseUrl string, log *zerolog.Logger) *Client {
	return &Client{baseUrl: baseUrl, log: log, ctx: ctx}
}

func (c *Client) InitGame(data battleships.GamePost) (battleships.Game, error) {
	method, endpoint := http.MethodPost, "/game"
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	_, headers, err := c.request(method, endpoint, "", body)
	if err != nil {
		return nil, err
	}

	game := NewGame(headers.Get(ApiTokenHeader), c.log)

	return game, nil
}

func (c *Client) UpdateBoard(game battleships.Game) error {
	method, endpoint := http.MethodGet, "/game/board"
	var body []byte

	res, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return err
	}

	var boardRes battleships.BoardGet
	if err = json.Unmarshal(res, &boardRes); err != nil {
		return err
	}

	game.SetBoard(boardRes.Board)
	return nil
}

func (c *Client) GameDesc(game battleships.Game) error {
	method, endpoint := http.MethodGet, "/game/desc"
	var body []byte

	res, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return err
	}

	var parsed battleships.GameGet
	if err = json.Unmarshal(res, &parsed); err != nil {
		return err
	}

	status := battleships.GameStatus{
		Status:     parsed.GameStatus,
		LastStatus: parsed.LastGameStatus,
		ShouldFire: parsed.ShouldFire,
		Timer:      parsed.Timer,
	}

	game.SetGameStatus(status)
	game.SetPlayer(NewPlayer(parsed.Nick, parsed.Desc))
	game.SetOpponent(NewPlayer(parsed.Opponent, parsed.OppDesc))
	return nil
}

func (c *Client) GameStatus(game battleships.Game) error {
	method, endpoint := http.MethodGet, "/game"
	var body []byte

	res, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return err
	}

	var parsed battleships.GameGet
	if err = json.Unmarshal(res, &parsed); err != nil {
		return err
	}

	status := battleships.GameStatus{
		Status:     parsed.GameStatus,
		LastStatus: parsed.LastGameStatus,
		ShouldFire: parsed.ShouldFire,
		Timer:      parsed.Timer,
	}

	game.SetGameStatus(status)
	return nil
}

func (c *Client) Fire(game battleships.Game, field string) (bool, error) {
	method, endpoint := http.MethodPost, "/game/fire"
	body, err := json.Marshal(struct {
		Field string `json:"coord"`
	}{field})
	if err != nil {
		return false, err
	}

	res, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return false, err
	}

	var parsed battleships.FireRes
	if err = json.Unmarshal(res, &parsed); err != nil {
		return false, err
	}

	err = c.GameStatus(game)
	if err != nil {
		return false, err
	}

	c.log.Debug().Interface("response", parsed).Msg("Fire response")

	return false, nil
}

func (c *Client) request(method, endpoint string, key string, body []byte) ([]byte, http.Header, error) {
	timeoutCtx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	reqUrl, err := url.JoinPath(c.baseUrl, endpoint)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(timeoutCtx, method, reqUrl, bytes.NewBuffer(body))
	if key != "" {
		req.Header.Set(ApiTokenHeader, key)
	}
	if err != nil {
		return nil, nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 {
		errMsg := fmt.Sprintf("Server returned code %d", res.StatusCode)

		var parsed map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&parsed); err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, nil, err
			}
		}

		if message, ok := parsed["message"]; ok {
			errMsg = fmt.Sprintf("Server returned code %d. Message: %v", res.StatusCode, message)
		}

		return nil, nil, errors.New(errMsg)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return resBody, res.Header, nil
}
