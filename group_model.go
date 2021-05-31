package hue

// Group struct that represents Philips Hue Group
//
// 0 (Zero) A special group containing all lights in the system, and is not returned by the ‘get all groups’ command. This group is not visible, and cannot be created, modified or deleted using the API.
type Group struct {
	ID     int
	Name   string      `json:"name"`   // A unique, editable name given to the group.
	Lights []string    `json:"lights"` // The IDs of the lights that are in the group.
	Type   string      `json:"type"`   // If not provided upon creation “LightGroup” is used. Can be “LightGroup”, “Room” or either “Luminaire” or “LightSource” if a Multisource Luminaire is present in the system.
	Action GroupAction `json:"action"` // The light state of one of the lamps in the group.
}

// GroupAction is used to execute actions on all lights in a group.
type GroupAction struct {
	On        bool      `json:"on"`     // On/Off state of the light. On=true, Off=false
	Bri       int       `json:"bri"`    // Brightness is a scale from 0 (the minimum the light is capable of) to 254 (the maximum). Note: a brightness of 0 is not off.e.g. “brightness”: 60 will set the light to a specific brightness.
	Hue       int       `json:"hue"`    // The hue value is a wrapping value between 0 and 65535. Both 0 and 65535 are red, 25500 is green and 46920 is blue.e.g. “hue”: 50000 will set the light to a specific hue.
	Sat       int       `json:"sat"`    // Saturation of the light. 254 is the most saturated (colored) and 0 is the least saturated (white).
	Effect    string    `json:"effect"` // The dynamic effect of the light, currently “none” and “colorloop” are supported. Other values will generate an error of type 7.Setting the effect to colorloop will cycle through all hues using the current brightness and saturation settings.
	XY        []float64 `json:"xy"`     // The x and y coordinates of a color in CIE color spaceThe first entry is the x coordinate and the second entry is the y coordinate. Both x and y must be between 0 and 1. If the specified coordinates are not in the CIE color space, the closest color to the coordinates will be chosen.
	Ct        int       `json:"ct"`     // The Mired Color temperature of the light. 2012 connected lights are capable of 153 (6500K) to 500 (2000K).
	Alert     string    `json:"alert"`
	Colormode string    `json:"colormode"`
}
