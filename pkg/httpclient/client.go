package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client interface of http client
type Client interface {
	Request(string, string, int, interface{}) error
}

type client struct {
	url    string
	client *http.Client
}

// New returns a http client
func New(url string) Client {
	return &client{
		url:    url,
		client: &http.Client{},
	}
}

func (c *client) Request(context, method string, statusCode int, body interface{}) error {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, c.url+context, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != statusCode {
		return fmt.Errorf(
			"xdpdropper request failed. Unexpected status code %d on %s. It should be %d",
			response.StatusCode,
			method,
			statusCode,
		)
	}

	return nil
}
