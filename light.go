package hue

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"

	funk "github.com/thoas/go-funk"
)

type SetStateParams struct {
	On             *bool      `json:"on,omitempty"`
	Bri            *uint8     `json:"bri,omitempty"`
	Hue            *uint16    `json:"hue,omitempty"`
	Sat            *uint8     `json:"sat,omitempty"`
	Effect         *string    `json:"effect,omitempty"`
	XY             *[]float64 `json:"xy,omitempty"`
	CT             *uint16    `json:"ct,omitempty"`
	Alert          *string    `json:"alert,omitempty"`
	TransitionTime *uint16    `json:"transitiontime,omitempty"`
	BriInc         *uint8     `json:"bri_inc,omitempty"`
	SatInc         *uint8     `json:"sat_inc,omitempty"`
	HueInc         *uint16    `json:"hue_inc,omitempty"`
	CtInc          *uint16    `json:"ct_inc,omitempty"`
	XYInc          *[]float32 `json:"xy_inc,omitempty"`
	Scene          *string    `json:"scene,omitempty"`
}

// LightService has functions for groups
type LightService service

const lightServiceName = "lights"

type renameRequest struct {
	Name string `json:"name"`
}

// GetAll returns a list of all lights that have been discovered by the bridge.
func (s *LightService) GetAll(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, path(lightServiceName, s.client.clientId), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(map[string]*Light)
	resp, err := s.client.do(ctx, req, parsed)
	if err != nil {
		return nil, resp, err
	}

	for i, l := range *parsed {
		id, _ := strconv.Atoi(i)
		l.ID = &id
	}

	result := funk.Values(*parsed).([]*Light)
	sort.Slice(result, func(i, j int) bool {
		return *result[i].ID < *result[j].ID
	})

	return result, resp, nil
}

// Get returns light by id
func (s *LightService) Get(ctx context.Context, id string) (*Light, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, path(lightServiceName, s.client.clientId, id), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(Light)
	resp, err := s.client.do(ctx, req, parsed)
	if err != nil {
		return nil, resp, err
	}

	return parsed, resp, nil
}

// GetNew returns a list of lights that were discovered the last time a search for new lights was performed.
func (s *LightService) GetNew(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.newRequest(http.MethodGet, path(lightServiceName, s.client.clientId, "new"), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(map[string]interface{})
	resp, err := s.client.do(ctx, req, parsed)
	if err != nil {
		return nil, resp, err
	}

	lights := []*Light{}
	for i, l := range *parsed {
		if i == "lastscan" {
			// Skip lastscan
		} else {
			light := &Light{}
			id, _ := strconv.Atoi(i)
			light.ID = &id
			light.Name = String(l.(map[string]interface{})["name"].(string))
			lights = append(lights, light)
		}
	}

	return lights, resp, nil
}

// Search starts searching for new lights
// The bridge will open the network for 40s.
func (s *LightService) Search(ctx context.Context) (bool, *Response, error) {
	req, err := s.client.newRequest(http.MethodPost, path(lightServiceName, s.client.clientId), nil)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 || (*apiResponses)[0].Error != nil {
		return false, resp, errors.New((*apiResponses)[0].Error.Description)
	}

	if s.client.verbose {
		s.client.logger.Info("%v", (*apiResponses)[0].Success)
	}

	return true, resp, nil
}

// Rename lights
func (s *LightService) Rename(ctx context.Context, id, name string) (bool, *Response, error) {
	payload := &renameRequest{name}
	req, err := s.client.newRequest(http.MethodPut, path(lightServiceName, s.client.clientId, id), payload)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 || (*apiResponses)[0].Error != nil {
		return false, resp, errors.New((*apiResponses)[0].Error.Description)
	}

	if s.client.verbose {
		s.client.logger.Info("%v", (*apiResponses)[0].Success)
	}

	return true, resp, nil
}

// SetState allows the user to turn the light on and off, modify the hue and effects.
func (s *LightService) SetState(ctx context.Context, id string, payload SetStateParams) ([]*ApiResponse, *Response, error) {
	req, err := s.client.newRequest(http.MethodPut, path(lightServiceName, s.client.clientId, id, "state"), payload)
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

// Delete a light from the bridge.
func (s *LightService) Delete(ctx context.Context, id string) (bool, *Response, error) {
	req, err := s.client.newRequest(http.MethodDelete, path(lightServiceName, s.client.clientId, id), nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.do(ctx, req, nil)
	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}
