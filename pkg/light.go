package hue

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	funk "github.com/thoas/go-funk"
)

// LightService has functions for groups
type LightService service

// RenameRequest
type RenameRequest struct {
	Name string `json:"name"`
}

type SetStateRequest struct {
	On             *bool      `json:"on,omitempty"`
	Bri            *uint8     `json:"bri,omitempty"`
	Hue            *uint16    `json:"hue,omitempty"`
	Sat            *uint8     `json:"sat,omitempty"`
	XY             *[]float32 `json:"xy,omitempty"`
	CT             *uint16    `json:"ct,omitempty"`
	Alert          *string    `json:"alert,omitempty"`
	Effect         *string    `json:"effect,omitempty"`
	TransitionTime *uint16    `json:"transitiontime,omitempty"`
	BriInc         *uint8     `json:"bri_inc,omitempty"`
	SatInc         *uint8     `json:"sat_inc,omitempty"`
	HueInc         *uint16    `json:"hue_inc,omitempty"`
	CTInc          *uint16    `json:"ct_inc,omitempty"`
	XYInc          *[]float32 `json:"xy_inc,omitempty"`
}

const lightServiceName = "lights"

// GetAll returns a list of all lights that have been discovered by the bridge.
func (s *LightService) GetAll(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, path(lightServiceName, s.client.Username), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(map[string]*Light)
	resp, err := s.client.Do(ctx, req, parsed)
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
	req, err := s.client.NewRequest(http.MethodGet, path(lightServiceName, s.client.Username, id), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(Light)
	resp, err := s.client.Do(ctx, req, parsed)
	if err != nil {
		return nil, resp, err
	}

	return parsed, resp, nil
}

// GetNew returns a list of lights that were discovered the last time a search for new lights was performed.
func (s *LightService) GetNew(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, path(lightServiceName, s.client.Username, "new"), nil)
	if err != nil {
		return nil, nil, err
	}

	parsed := new(map[string]interface{})
	resp, err := s.client.Do(ctx, req, parsed)
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
	req, err := s.client.NewRequest(http.MethodPost, path(lightServiceName, s.client.Username), nil)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.Do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if len(*apiResponses) > 0 && len((*apiResponses)[0].Success) != 0 {
		if s.client.Verbose {
			s.client.logger.Info("%v", (*apiResponses)[0].Success["/lights"])
		}
		return true, resp, nil
	}

	return false, resp, nil
}

// Rename lights
func (s *LightService) Rename(ctx context.Context, id, name string) (bool, *Response, error) {
	payload := &RenameRequest{name}
	req, err := s.client.NewRequest(http.MethodPut, path(lightServiceName, s.client.Username, id), payload)
	if err != nil {
		return false, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.Do(ctx, req, apiResponses)
	if err != nil {
		return false, resp, err
	}

	if len(*apiResponses) > 0 && len((*apiResponses)[0].Success) != 0 {
		if s.client.Verbose {
			s.client.logger.Info("%v", (*apiResponses)[0].Success[fmt.Sprintf("/lights/%s/name", id)])
		}
		return true, resp, nil
	}

	return false, resp, nil
}

// SetState allows the user to turn the light on and off, modify the hue and effects.
func (s *LightService) SetState(ctx context.Context, id string, payload *SetStateRequest) ([]*ApiResponse, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPut, path(lightServiceName, s.client.Username, id, "state"), payload)
	if err != nil {
		return nil, nil, err
	}

	apiResponses := new([]*ApiResponse)
	resp, err := s.client.Do(ctx, req, apiResponses)
	if err != nil {
		return nil, resp, err
	}

	return *apiResponses, resp, nil
}

// Delete a light from the bridge.
func (s *LightService) Delete(ctx context.Context, id string) (bool, *Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, path(lightServiceName, s.client.Username, id), nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}
