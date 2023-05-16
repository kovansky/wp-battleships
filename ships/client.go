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

func (c *Client) Abandon(game battleships.Game) error {
	method, endpoint := http.MethodGet, "/game/abandon"
	var body []byte

	_, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return err
	}

	return nil
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

	boardFields := make(map[string]battleships.FieldState, len(boardRes.Board))
	for _, field := range boardRes.Board {
		boardFields[field] = battleships.FieldStateShip
	}

	game.SetBoard(boardFields)
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

	board := game.Board()

	if board != nil {
		for _, shot := range parsed.OppShots {
			if _, ok := board[shot]; ok &&
				(board[shot] == battleships.FieldStateShip ||
					board[shot] == battleships.FieldStateHit) {
				board[shot] = battleships.FieldStateHit
			} else {
				board[shot] = battleships.FieldStateMiss
			}
		}
		game.SetBoard(board)
	}

	return nil
}

func (c *Client) Refresh(game battleships.Game) error {
	method, endpoint := http.MethodGet, "/game/refresh"
	var body []byte

	_, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) PlayerStats(nick string) (battleships.PlayerStats, error) {
	method := http.MethodGet
	endpoint, err := url.JoinPath("/game/stats", nick)
	if err != nil {
		return battleships.PlayerStats{}, err
	}

	var body []byte

	res, _, err := c.request(method, endpoint, "", body)
	if err != nil {
		return battleships.PlayerStats{}, err
	}

	var parsed struct {
		Stats battleships.PlayerStats `json:"stats"`
	}
	if err = json.Unmarshal(res, &parsed); err != nil {
		return battleships.PlayerStats{}, err
	}

	return parsed.Stats, nil
}

func (c *Client) ListPlayers() ([]battleships.Player, error) {
	var players []battleships.Player

	method, endpoint := http.MethodGet, "/game/list"
	var body []byte

	res, _, err := c.request(method, endpoint, "", body)
	if err != nil {
		return players, err
	}

	type responseType struct {
		Nick string `json:"nick"`
	}
	var parsed []struct {
		Nick string `json:"nick"`
	}
	if err = json.Unmarshal(res, &parsed); err != nil {
		return players, err
	}

	parsed = append(parsed, responseType{Nick: "WP_Bot"})

	for _, player := range parsed {
		player, err := c.PlayerStats(player.Nick)
		if err != nil {
			continue
		}

		players = append(players, NewPlayerFromStats(player))
	}

	return players, nil
}

func (c *Client) Fire(game battleships.Game, field string) (battleships.ShotState, error) {
	method, endpoint := http.MethodPost, "/game/fire"
	body, err := json.Marshal(struct {
		Field string `json:"coord"`
	}{field})
	if err != nil {
		return "", err
	}

	res, _, err := c.request(method, endpoint, game.Key(), body)
	if err != nil {
		return "", err
	}

	var parsed battleships.FireRes
	if err = json.Unmarshal(res, &parsed); err != nil {
		return "", err
	}

	err = c.GameStatus(game)
	if err != nil {
		return "", err
	}

	return parsed.Result, nil
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
