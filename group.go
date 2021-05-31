package hue

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"

	funk "github.com/thoas/go-funk"
)

// GroupService has functions for groups
type GroupService service

const groupServiceName = "groups"

const groupTypeLight = "LightGroup"
const groupTypeRoom = "Room"

type createGroupRequest struct {
	Lights []string `json:"lights,omitempty"`
	Name   string   `json:"name,omitempty"`
	Type   string   `json:"type,omitempty"`
	Class  *string  `json:"class,omitempty"`
}

type updateGroupRequest struct {
	Lights []string `json:"lights,omitempty"`
	Name   *string  `json:"name,omitempty"`
	Class  *string  `json:"class,omitempty"`
}

func (s *GroupService) groupServicePath(params ...string) string {
	return s.client.path(groupServiceName, params...)
}

// GetAll returns all groups
func (s *GroupService) GetAll(ctx context.Context) ([]*Group, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, s.groupServicePath(), nil)
	if err != nil {
		return nil, nil, err
	}

	groups := make(map[string]*Group)
	resp, err := s.client.do(ctx, req, &groups)
	if err != nil {
		return nil, resp, err
	}

	for i, g := range groups {
		id, _ := strconv.Atoi(i)
		g.ID = id
	}

	orderedGroups := funk.Values(groups).([]*Group)
	sort.Slice(orderedGroups, func(i, j int) bool {
		return orderedGroups[i].ID < orderedGroups[j].ID
	})

	return orderedGroups, resp, nil
}

// CreateGroup creates light group and returns id of the created group
func (s *GroupService) CreateGroup(ctx context.Context, name string, lights []string) (string, *Response, error) {
	payload := &createGroupRequest{
		Name:   name,
		Lights: lights,
		Type:   groupTypeLight,
	}
	req, err := s.client.newRequest(http.MethodPost, s.groupServicePath(), payload)
	if err != nil {
		return "", nil, err
	}

	apiResponses := new([]ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return "", resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 || (*apiResponses)[0].Error != nil {
		return "", resp, errors.New((*apiResponses)[0].Error.Description)
	}

	// Get first success message
	successResponse := (*apiResponses)[0].Success

	return successResponse["id"].(string), resp, nil
}

// CreateGroup creates light room and returns id of the room
func (s *GroupService) CreateRoom(ctx context.Context, name string, lights []string) (string, *Response, error) {
	payload := &createGroupRequest{
		Name:   name,
		Lights: lights,
		Type:   groupTypeRoom,
		Class:  String(name),
	}
	req, err := s.client.newRequest(http.MethodPost, s.groupServicePath(), payload)
	if err != nil {
		return "", nil, err
	}

	apiResponses := new([]ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return "", resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 || (*apiResponses)[0].Error != nil {
		return "", resp, errors.New((*apiResponses)[0].Error.Description)
	}

	// Get first success message
	successResponse := (*apiResponses)[0].Success

	return successResponse["id"].(string), resp, nil
}

// Get returns the group by id
func (s *GroupService) Get(ctx context.Context, id string) (*Group, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, s.groupServicePath(id), nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(Group)
	resp, err := s.client.do(ctx, req, group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// Update updates group by id
func (s *GroupService) Update(ctx context.Context, id string, name *string, lights []string, class *string) (bool, *Response, error) {
	payload := &updateGroupRequest{
		Name:   name,
		Lights: lights,
		Class:  class,
	}
	req, err := s.client.newRequest(http.MethodPut, s.groupServicePath(id), payload)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 {
		return false, resp, errors.New("the bridge didn't return valid response")
	}

	// If response has any error, return as fail
	for _, r := range *apiResponses {
		if r.Error != nil {
			return false, resp, errors.New(r.Error.Description)
		}
	}

	return true, resp, nil
}

// SetState updates state of the group
func (s *GroupService) SetState(ctx context.Context, id string, payload SetStateParams) ([]*ApiResponse, *Response, error) {
	req, err := s.client.newRequest(http.MethodPut, s.groupServicePath(id, "action"), payload)
	if err != nil {
		return nil, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return nil, resp, err
	}

	return *apiResponses, resp, nil
}

// Delete removes the group
func (s *GroupService) Delete(ctx context.Context, id string) (bool, *Response, error) {
	req, err := s.client.newRequest(http.MethodDelete, s.groupServicePath(id), nil)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]map[string]string)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 {
		return false, resp, errors.New("the bridge didn't return valid response")
	}

	return true, resp, nil
}
