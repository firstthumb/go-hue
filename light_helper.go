package hue

import (
	"context"
)

func (s *LightService) TurnOn(ctx context.Context, id string) ([]*ApiResponse, *Response, error) {
	return s.SetState(ctx, id, SetStateParams{On: Bool(true)})
}

func (s *LightService) TurnOff(ctx context.Context, id string) ([]*ApiResponse, *Response, error) {
	return s.SetState(ctx, id, SetStateParams{On: Bool(false)})
}

func (s *LightService) TurnOnAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		s.SetState(ctx, id, SetStateParams{On: Bool(true)})
	}
}

func (s *LightService) TurnOffAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		s.SetState(ctx, id, SetStateParams{On: Bool(false)})
	}
}
