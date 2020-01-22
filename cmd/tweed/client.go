package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type client struct {
	http     *http.Client
	url      string
	username string
	password string
}

func Connect(url, username, password string) *client {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return &client{
		url:      strings.TrimSuffix(url, "/"),
		username: username,
		password: password,
		http:     http.DefaultClient,
	}
}

func (c *client) do(req *http.Request, out interface{}) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	if Tweed.Debug {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUG: @W{unable to dump request:} @R{%s}\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n\n", string(b))
		}
	}
	res, err := c.http.Do(req)
	if res != nil && Tweed.Debug {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUG: @W{unable to dump response:} @R{%s}\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n\n", string(b))
		}
	}
	if err != nil || out == nil {
		return res, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}

	var e api.ErrorResponse
	if res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 404 || res.StatusCode == 500 {
		if err = json.Unmarshal(b, &e); err != nil {
			return res, err
		}
		return res, e
	}

	return res, json.Unmarshal(b, &out)
}

func (c *client) request(method, path string, in interface{}) (*http.Request, error) {
	var body bytes.Buffer
	if in != nil {
		b, err := json.MarshalIndent(in, "", "  ")
		if err != nil {
			return nil, err
		}
		body.Write(b)
	}

	path = strings.TrimPrefix(path, "/")
	req, err := http.NewRequest(method, c.url+"/"+path, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	if in != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (c *client) GET(path string, out interface{}) error {
	req, err := c.request("GET", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) POST(path string, in, out interface{}) error {
	req, err := c.request("POST", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) PUT(path string, in, out interface{}) error {
	req, err := c.request("PUT", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) DELETE(path string, out interface{}) error {
	req, err := c.request("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}
