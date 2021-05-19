package hue

import (
	"context"
	"net/http"
	"sort"
	"strconv"

	funk "github.com/thoas/go-funk"
)

// GroupService has functions for groups
type GroupService service

const groupServiceName = "groups"

// GetAll returns all groups
func (s *GroupService) GetAll(ctx context.Context) ([]*Group, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, path(groupServiceName, s.client.clientId), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(map[string]*Group)
	resp, err := s.client.do(ctx, req, parsed)
	if err != nil {
		return nil, resp, err
	}

	for i, g := range *parsed {
		id, _ := strconv.Atoi(i)
		g.ID = &id
	}

	result := funk.Values(*parsed).([]*Group)
	sort.Slice(result, func(i, j int) bool {
		return *result[i].ID < *result[j].ID
	})

	return result, resp, nil
}
