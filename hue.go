package hue

import (
	hue "github.com/firstthumb/go-hue/pkg"
)

func Discover() (*hue.Client, error) {
	return hue.Discover()
}

func NewClient(host, username string, opts *hue.ClientOptions) *hue.Client {
	return hue.NewClient(host, username, opts)
}
