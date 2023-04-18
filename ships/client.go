package ships

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

const ApiTokenHeader = "X-Auth-Token"

type Client struct {
	baseUrl string
	client  http.Client
	log     *zerolog.Logger
}

func NewClient(baseUrl string, log *zerolog.Logger) *Client {
	return &Client{baseUrl: baseUrl, log: log}
}

func (c *Client) InitGame(withBot bool) (*Game, error) {
	method, endpoint := http.MethodPost, "/shps"
	body, err := json.Marshal(struct {
		Wpbot bool `json:"wpbot"`
	}{
		withBot,
	})
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

	url := fmt.Sprintf("%s%s", c.baseUrl, endpoint)

	req, err := http.NewRequestWithContext(timeoutCtx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, nil, errors.New(fmt.Sprintf("Server returned code %d", res.StatusCode))
	}

	var parsed map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, nil, err
		}
	}

	return parsed, res.Header, nil
}
