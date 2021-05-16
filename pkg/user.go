package hue

import (
	"context"
	"net/http"
)

type UserService service

type User struct {
	Lights map[string]Light `json:"lights,omitempty"`
}

func (s *UserService) Login(ctx context.Context, username string) (*User, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, username, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(User)
	resp, err := s.client.Do(ctx, req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, nil
}
