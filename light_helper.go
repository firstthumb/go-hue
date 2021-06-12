package hue

import (
	"context"
	"errors"
	"image/color"
)

// TurnOn sets on status as true
func (s *LightService) TurnOn(ctx context.Context, id string) error {
	_, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(true)})
	return err
}

// TurnOff sets on status as false
func (s *LightService) TurnOff(ctx context.Context, id string) error {
	_, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(false)})
	return err
}

// TurnOnAll sets on status as true
func (s *LightService) TurnOnAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		apiResponses, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(true)})
		if err != nil {
			s.client.logger.Info("Turning on failed", "ApiResponses", apiResponses)
		}
	}
}

// TurnOffAll sets on status as false
func (s *LightService) TurnOffAll(ctx context.Context, ids ...string) {
	for _, id := range ids {
		apiResponses, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(false)})
		if err != nil {
			s.client.logger.Info("Turning off failed", "ApiResponses", apiResponses)
		}
	}
}

// SetColor changes the color of lamp with color
func (s *LightService) SetColor(ctx context.Context, id string, clr color.Color) error {
	hueVal, satVal, briVal, ok := colorToHSV(clr)
	if !ok {
		return errors.New("the color is not supported")
	}

	apiResponses, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(true), Hue: &hueVal, Sat: &satVal, Bri: &briVal})
	if err != nil {
		return err
	}

	s.client.logger.Info("Set state successful", "ApiResponses", apiResponses)

	return nil
}

// SetColorHex changes the color of lamp with hex color code
func (s *LightService) SetColorHex(ctx context.Context, id string, hex string) error {
	hueVal, satVal, briVal, ok := hexColorToHSV(hex)
	if !ok {
		return errors.New("the color is not supported")
	}

	apiResponses, _, err := s.SetState(ctx, id, SetStateParams{On: Bool(true), Hue: &hueVal, Sat: &satVal, Bri: &briVal})
	if err != nil {
		return err
	}

	s.client.logger.Info("Set state successful", "ApiResponses", apiResponses)

	return nil
}
