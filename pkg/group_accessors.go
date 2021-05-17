package hue

// GetName returns human readable name of the group.
// If name is not specified one is generated for you (default name is “Group”)
func (g *Group) GetName() string {
	if g == nil || g.Name == nil {
		return ""
	}
	return *g.Name
}

// GetLights returns the ordered set of light ids from the lights which are in the group.
// This resource shall contain an array of at least one element with the exception of the “Room” type: The Room type may contain an empty lights array.
// Each element can appear only once. Order of lights on creation is preserved. A light id must be an existing light resource in /lights.
// If an invalid lights resource is given, error 7 shall be returned and the group is not created. There shall be no change in the lights.
// Light id can be null if a group has been automatically create by the bridge and a light source is not yet available
func (g *Group) GetLights() []string {
	if g == nil || g.Lights == nil {
		return nil
	}

	r := make([]string, len(g.Lights))
	for _, v := range g.Lights {
		r = append(r, *v)
	}

	return r
}

// GetType returns type of the Group. If not provided on creation a “LightGroup” is created. Supported types:
// LightGroup	1.4		Default
// Luminaire	1.4		multisource luminaire
// LightSource	1.4		multisource luminaire
// Room			1.11	Represents a room
// Entertainment	1.22	Represents an entertainment setup
// Zone			1.30	Represents a zone
func (g *Group) GetType() string {
	if g == nil || g.Name == nil {
		return ""
	}
	return *g.Type
}

// IsOn returns On/Off state of the light. On=true, Off=false
func (g *Group) IsOn() bool {
	if g == nil || g.Action == nil || g.Action.On == nil {
		return false
	}
	return *g.Action.On
}

// GetBrightness returns which is a scale from 0 (the minimum the light is capable of) to 254 (the maximum).
// Note: a brightness of 0 is not off.e.g. “brightness”: 60 will set the light to a specific brightness.
func (g *Group) GetBrightness() int {
	if g == nil || g.Action == nil || g.Action.Bri == nil {
		return 0
	}
	return *g.Action.Bri
}

// TODO: Add rest properties
