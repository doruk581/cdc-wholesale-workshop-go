package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/doruk581/cdc-wholesale-workshop-go/model"
	"io"
	"net/http"
	"net/url"
)

// Client is our consumer interface to the Listing API
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
	Token      string
}

// WithToken applies a token to the
func (c *Client) WithToken(token string) *Client {
	c.Token = token
	return c
}

// GetProduct gets a single product from the API
func (c *Client) GetProduct(id int) (*model.Product, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/product/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var product model.Product
	res, err := c.do(req, &product)

	if res != nil {
		switch res.StatusCode {
		case http.StatusNotFound:
			return nil, ErrNotFound
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		}
	}

	if err != nil {
		return nil, ErrUnavailable
	}

	return &product, err

}

// GetProducts gets all products from the API
func (c *Client) GetProducts() ([]model.Product, error) {
	req, err := c.newRequest("GET", "/products", nil)
	if err != nil {
		return nil, err
	}
	var products []model.Product
	_, err = c.do(req, &products)

	return products, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Admin Service")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

var (
	ErrNotFound = errors.New("not found")

	ErrUnauthorized = errors.New("unauthorized")

	ErrUnavailable = errors.New("api unavailable")
)
