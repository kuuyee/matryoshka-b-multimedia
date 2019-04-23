package multimedia

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var defaultClient = new(http.Client)

// Client is multimedia client struct
type Client struct {
	Client  *http.Client
	BaseURL string
}

func (c Client) getClient() *http.Client {
	if c.Client != nil {
		return c.Client
	}
	return defaultClient
}

func (c Client) buildRequest(method string, partialURL string, body io.Reader) *http.Request {
	fullURL := strings.TrimSuffix(c.BaseURL, "/") + "/" + strings.TrimPrefix(partialURL, "/")
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		log.Panicf("error while building request %s %s: %v", method, fullURL, err)
	}
	return req
}

func (c Client) makeJSONRequest(method string, partialURL string, reqData interface{}, respData interface{}) (int, error) {
	var body io.Reader
	if reqData != nil {
		bodyJSON, err := json.Marshal(reqData)
		if err != nil {
			return 0, err
		}
		body = bytes.NewBuffer(bodyJSON)
	}
	req := c.buildRequest(method, partialURL, body)
	if reqData != nil {
		req.Header.Set("content-type", "application/json")
	}
	resp, err := c.getClient().Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return resp.StatusCode, HTTPCodeError{resp}
	}
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, json.Unmarshal(respJSON, respData)
}
