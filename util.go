package hue

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func Bool(v bool) *bool { return &v }

func Int(v int) *int { return &v }

func UInt8(v uint8) *uint8 { return &v }

func Int64(v int64) *int64 { return &v }

func String(v string) *string { return &v }

func Slice(v []string) *[]string { return &v }

func colorToHSV(clr color.Color) (uint16, uint8, uint8, bool) {
	hueclr, ok := colorful.MakeColor(clr)
	if !ok {
		return 0, 0, 0, false
	}

	hue, sat, val := hueclr.Hsv()

	hueVal := uint16((hue * 65535) / 360)
	satVal := uint8(sat * 255)
	briVal := uint8(val * 255)

	return hueVal, satVal, briVal, true
}

func hexColorToHSV(hexColor string) (uint16, uint8, uint8, bool) {
	hueclr, err := colorful.Hex(hexColor)
	if err != nil {
		return 0, 0, 0, false
	}

	hue, sat, val := hueclr.Hsv()

	hueVal := uint16((hue * 65535) / 360)
	satVal := uint8(sat * 255)
	briVal := uint8(val * 255)

	return hueVal, satVal, briVal, true
}
