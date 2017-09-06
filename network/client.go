package network

import "net/url"

type Client struct {
}

func NewClient(u *url.URL) *Client {
	return &Client{}
}
