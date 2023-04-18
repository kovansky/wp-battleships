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
}

func NewClient(baseUrl string, log *zerolog.Logger) *Client {
	return &Client{baseUrl: baseUrl, log: log}
}

func (c *Client) InitGame(data battleships.GamePost) (battleships.Game, error) {
	method, endpoint := http.MethodPost, "/game"
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	_, headers, err := c.request(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	game := NewGame(headers.Get(ApiTokenHeader), c.log)

	return game, nil
}

func (c *Client) request(method, endpoint string, body []byte) (map[string]interface{}, http.Header, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqUrl, err := url.JoinPath(c.baseUrl, endpoint)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(timeoutCtx, method, reqUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var parsed map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, nil, err
		}
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		errMsg := fmt.Sprintf("Server returned code %d", res.StatusCode)

		if message, ok := parsed["message"]; ok {
			errMsg = fmt.Sprintf("Server returned code %d. Message: %v", res.StatusCode, message)
		}

		return nil, nil, errors.New(errMsg)
	}

	return parsed, res.Header, nil
}
