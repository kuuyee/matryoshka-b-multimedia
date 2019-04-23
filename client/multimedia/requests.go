package multimedia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"

	"github.com/kuuyee/matryoshka-b-multimedia/model"
)

// FileMeta fetches file metadata
func (c Client) FileMeta(service Service, ident string) (*model.Meta, error) {
	res := new(model.Meta)
	if _, err := c.makeJSONRequest("GET", fmt.Sprintf("/%s/%s/meta", service, ident), nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// FetchFile fetches file data
func (c Client) FetchFile(service Service, ident string, params map[string]string) (io.ReadCloser, error) {
	urlBuilder := bytes.NewBufferString(fmt.Sprintf("/%s/%s", service, ident))
	if params != nil {
		urlBuilder.WriteByte('?')
		for k, v := range params {
			urlBuilder.WriteString(url.QueryEscape(k))
			urlBuilder.WriteByte('=')
			urlBuilder.WriteString(url.QueryEscape(v))
			urlBuilder.WriteByte('&')
		}
		urlBuilder.Truncate(urlBuilder.Len() - 1)
	}
	req := c.buildRequest("GET", urlBuilder.String(), nil)
	resp, err := c.getClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, HTTPCodeError{resp}
	}
	return resp.Body, nil
}

// PostFile posts a new multimedia file
func (c Client) PostFile(service Service, fileName string, data io.Reader, params map[string]string) (*model.Meta, error) {
	urlBuilder := bytes.NewBufferString(fmt.Sprintf("/%s", service))
	if params != nil {
		urlBuilder.WriteByte('?')
		for k, v := range params {
			urlBuilder.WriteString(url.QueryEscape(k))
			urlBuilder.WriteByte('=')
			urlBuilder.WriteString(url.QueryEscape(v))
			urlBuilder.WriteByte('&')
		}
		urlBuilder.Truncate(urlBuilder.Len() - 1)
	}

	bodyWriter := bytes.NewBuffer([]byte{})
	postForm := multipart.NewWriter(bodyWriter)
	fileWriter, err := postForm.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fileWriter, data); err != nil {
		return nil, err
	}
	postForm.Close()
	req := c.buildRequest("POST", urlBuilder.String(), bodyWriter)
	req.Header.Set("content-type", postForm.FormDataContentType())

	resp, err := c.getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, HTTPCodeError{resp}
	}
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := new(model.Meta)
	if err := json.Unmarshal(respJSON, res); err != nil {
		return nil, err
	}
	return res, nil
}
