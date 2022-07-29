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

	if out == nil {
		// nothing is expected, make sure its a 200 resp code
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(data) == `{"data":null}` {
			return nil
		}
		if string(data) == "null" {
			return nil
		}
		if string(data) == "" {
			return nil
		}
		// its a json that represnets an error, just reutrn it
		return fmt.Errorf("json failed to decode post message: '%s'", string(data))
	}
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

func (c *Client) decodeResp(resp *http.Response, out interface{}) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.logger.Printf("[TRACE] Http response: data, %s", string(data))

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
