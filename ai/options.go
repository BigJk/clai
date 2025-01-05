package ai

import "net/http"

type Options func(*Client)

// WithURL sets the api url
func WithURL(url string) Options {
	return func(c *Client) {
		c.URL = url
	}
}

// WithAPIKey sets the api key
func WithAPIKey(apiKey string) Options {
	return func(c *Client) {
		c.APIKey = apiKey
	}
}

// WithClient sets the http client
func WithClient(client *http.Client) Options {
	return func(c *Client) {
		c.client = client
	}
}

// WithModel sets the model
func WithModel(model string) Options {
	return func(c *Client) {
		c.Model = model
	}
}

// WithOpenAI sets the openai url
func WithOpenAI() Options {
	return func(c *Client) {
		c.URL = "https://api.openai.com/v1/chat/completions"
	}
}

// WithOpenRouter sets the openrouter url
func WithOpenRouter() Options {
	return func(c *Client) {
		c.URL = "https://openrouter.ai/api/v1/chat/completions"
	}
}
