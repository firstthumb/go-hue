package hue

import "context"

func (s *GroupService) TurnOn(ctx context.Context, id string) ([]*ApiResponse, *Response, error) {
	return s.SetState(ctx, id, SetStateParams{On: Bool(true)})
}

func (s *GroupService) TurnOff(ctx context.Context, id string) ([]*ApiResponse, *Response, error) {
	return s.SetState(ctx, id, SetStateParams{On: Bool(false)})
}

func (s *GroupService) TurnOnAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		s.SetState(ctx, id, SetStateParams{On: Bool(true)})
	}
}

func (s *GroupService) TurnOffAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		s.SetState(ctx, id, SetStateParams{On: Bool(false)})
	}
}
