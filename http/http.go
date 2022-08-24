package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// https://ethereum.github.io/beacon-APIs/#/

type Client struct {
	url    string
	logger *log.Logger
}

func New(url string) *Client {
	return &Client{url: url, logger: log.New(ioutil.Discard, "", 0)}
}

func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
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

func (c *Client) Get(path string, out interface{}) error {
	resp, err := http.Get(c.url + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.logger.Printf("[TRACE] Get request: path, %s", path)

	if err := c.decodeResp(resp, out); err != nil {
		return err
	}
	return nil
}

var (
	ErrorBadRequest          = fmt.Errorf("bad request (400)")
	ErrorNotFound            = fmt.Errorf("not found (404)")
	ErrorInternalServerError = fmt.Errorf("internal server error (500)")
	ErrorServiceUnavailable  = fmt.Errorf("service unavailable (503)")
)

func (c *Client) decodeResp(resp *http.Response, out interface{}) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	errorFn := func(code error) error {
		return fmt.Errorf("status code != 200: %s %w", string(data), code)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest: // 400
		return errorFn(ErrorBadRequest)
	case http.StatusNotFound: // 404
		return errorFn(ErrorNotFound)
	case http.StatusInternalServerError: // 500
		return errorFn(ErrorInternalServerError)
	case http.StatusServiceUnavailable: // 503
		return errorFn(ErrorServiceUnavailable)
	default:
		return errorFn(fmt.Errorf("%d", resp.StatusCode))
	}

	c.logger.Printf("[TRACE] Http response: data, %s", string(data))

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
	if err := Unmarshal(output.Data, &out); err != nil {
		return err
	}
	return nil
}
