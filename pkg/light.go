package hue

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	funk "github.com/thoas/go-funk"
)

// Service
type LightService service

// Request
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

func path(params ...string) string {
	if len(params) == 1 {
		return fmt.Sprintf("%v/lights", params[0])
	} else if len(params) == 2 {
		return fmt.Sprintf("%v/lights/%v", params[0], params[1])
	}

	return fmt.Sprintf("%v/lights/%v/%v", params[0], params[1], params[2])
}

func (s *LightService) GetAll(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, path(s.client.Username), nil)
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

func (s *LightService) Get(ctx context.Context, id string) (*Light, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, path(s.client.Username, id), nil)
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

func (s *LightService) GetNew(ctx context.Context) ([]*Light, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, path(s.client.Username, "new"), nil)
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

func (s *LightService) Search(ctx context.Context) (bool, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, path(s.client.Username), nil)
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

func (s *LightService) Rename(ctx context.Context, id, name string) (bool, *Response, error) {
	payload := &RenameRequest{name}
	req, err := s.client.NewRequest(http.MethodPut, path(s.client.Username, id), payload)
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

func (s *LightService) SetState(ctx context.Context, id string, payload *SetStateRequest) ([]*ApiResponse, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPut, path(s.client.Username, id, "state"), payload)
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

func (s *LightService) Delete(ctx context.Context, id string) (bool, *Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, path(s.client.Username, id), nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}
