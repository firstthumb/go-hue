package hue

import "context"

func (s *GroupService) TurnOn(ctx context.Context, id string) error {
	_, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(true)})
	return err
}

func (s *GroupService) TurnOff(ctx context.Context, id string) error {
	_, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(false)})
	return err
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
