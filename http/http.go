package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// https://ethereum.github.io/beacon-APIs/#/

type Config struct {
	logger        *log.Logger
	untrackedKeys bool
}

type ConfigOption func(*Config)

func WithLogger(logger *log.Logger) ConfigOption {
	return func(c *Config) {
		c.logger = logger
	}
}

func WithUntrackedKeys() ConfigOption {
	return func(c *Config) {
		c.untrackedKeys = true
	}
}

type Client struct {
	url    string
	config *Config
}

func New(url string, opts ...ConfigOption) *Client {
	config := &Config{
		logger: log.New(io.Discard, "", 0),
	}
	for _, opt := range opts {
		opt(config)
	}

	return &Client{url: url, config: config}
}

func (c *Client) SetLogger(logger *log.Logger) {
	c.config.logger = logger
}

func (c *Client) Post(path string, input interface{}, out interface{}) error {
	postBody, err := Marshal(input)
	if err != nil {
		return err
	}
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(c.url+path, "application/json", responseBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.decodeResp(resp, out); err != nil {
		return err
	}
	return nil
}

func (c *Client) Status(path string) (bool, error) {
	resp, err := http.Get(c.url + path)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	errorMsg, ok := httpErrorMapping[resp.StatusCode]
	if ok {
		return false, errorMsg
	}
	return false, fmt.Errorf("status code != 200: %d", resp.StatusCode)
}

func (c *Client) Get(path string, out interface{}) error {
	resp, err := http.Get(c.url + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.config.logger.Printf("[TRACE] Get request: path, %s", path)

	if err := c.decodeResp(resp, out); err != nil {
		return err
	}
	return nil
}

var (
	ErrorIncompleteData      = fmt.Errorf("incomplete data (206)")
	ErrorBadRequest          = fmt.Errorf("bad request (400)")
	ErrorNotFound            = fmt.Errorf("not found (404)")
	ErrorInternalServerError = fmt.Errorf("internal server error (500)")
	ErrorServiceUnavailable  = fmt.Errorf("service unavailable (503)")
)

var httpErrorMapping = map[int]error{
	http.StatusPartialContent:      ErrorIncompleteData,
	http.StatusBadRequest:          ErrorBadRequest,
	http.StatusNotFound:            ErrorNotFound,
	http.StatusInternalServerError: ErrorInternalServerError,
	http.StatusServiceUnavailable:  ErrorServiceUnavailable,
}

func (c *Client) decodeResp(resp *http.Response, out interface{}) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg, ok := httpErrorMapping[resp.StatusCode]
		if ok {
			return errorMsg
		}
		return fmt.Errorf("status code != 200: %s %d", string(data), resp.StatusCode)
	}

	c.config.logger.Printf("[TRACE] Http response: data, %s", string(data))

	if resp.Request.Method == http.MethodPost && out == nil {
		// post methods that expects no output
		if string(data) == `{"data":null}` {
			return nil
		}
		if string(data) == "null" {
			return nil
		}
		if string(data) == "" {
			return nil
		}
		return fmt.Errorf("json failed to decode post message: '%s'", string(data))
	}

	var output struct {
		Data json.RawMessage `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &output); err != nil {
		return err
	}
	if err := Unmarshal(output.Data, &out, c.config.untrackedKeys); err != nil {
		return err
	}
	return nil
}
