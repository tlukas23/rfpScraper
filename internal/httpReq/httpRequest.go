package httpReq

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func MakeHttpRequest(method string, urlEndpoint string, bodyFields map[string]string, headers map[string]string) ([]byte, error) {
	client := http.Client{}
	data := url.Values{}
	for key, value := range bodyFields {
		data.Set(key, value)
	}
	req, err := http.NewRequest(method, urlEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	return body, nil
}

func MakeHttpQuery(method string, urlEndpoint string, params map[string]string, headers map[string]string) ([]byte, error) {
	client := http.Client{}

	queryUrl, err := url.Parse(urlEndpoint)
	if err != nil {
		return nil, err
	}
	values := queryUrl.Query()
	for key, value := range params {
		values.Add(key, value)
	}
	queryUrl.RawQuery = values.Encode()

	req, err := http.NewRequest(method, queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	response.Body.Close()
	return body, nil
}
