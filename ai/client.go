package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	URL    string
	APIKey string
	Model  string

	client *http.Client
}

func NewClient(opts ...Options) *Client {
	c := &Client{
		URL:    "https://api.openai.com/v1/chat/completions",
		APIKey: "",
		client: &http.Client{
			Timeout: time.Second * 60 * 5,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Do(messages []Message) (string, error) {
	req := Request{
		Model:    c.Model,
		Messages: messages,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(body))
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", errors.New("no response")
	}

	return res.Choices[0].Message.Content, nil
}
