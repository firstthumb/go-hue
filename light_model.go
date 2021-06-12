package hue

type Light struct {
	ID               int          `json:"-"`
	State            State        `json:"state,omitempty"`
	SWUpdate         SWUpdate     `json:"swupdate,omitempty"`
	Type             string       `json:"type,omitempty"`
	Name             string       `json:"name"`
	ModelId          string       `json:"modelid,omitempty"`
	ManufacturerName string       `json:"manufacturername,omitempty"`
	ProductName      string       `json:"productname,omitempty"`
	Capabilities     Capabilities `json:"capabilities,omitempty"`
	Config           Config       `json:"config,omitempty"`
	UniqueId         string       `json:"uniqueid,omitempty"`
	SWVersion        string       `json:"swversion,omitempty"`
	SWConfigId       string       `json:"swconfigid,omitempty"`
	ProductId        string       `json:"productid,omitempty"`
}

type Config struct {
	Archetype string  `json:"archetype"`
	Function  string  `json:"function"`
	Direction string  `json:"direction"`
	Startup   Startup `json:"startup"`
}

type Startup struct {
	Mode       string `json:"mode"`
	Configured bool   `json:"configured"`
}

type Capabilities struct {
	Certified bool      `json:"certified"`
	Control   Control   `json:"control"`
	Streaming Streaming `json:"streaming"`
}

type Streaming struct {
	Renderer bool `json:"renderer"`
	Proxy    bool `json:"proxy"`
}

type Control struct {
	Mindimlevel    int         `json:"mindimlevel"`
	Maxlumen       int         `json:"maxlumen"`
	Colorgamuttype string      `json:"colorgamuttype"`
	Colorgamut     [][]float64 `json:"colorgamut"`
	Ct             Ct          `json:"ct"`
}

type Ct struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SWUpdate struct {
	State       string `json:"state"`
	Lastinstall string `json:"lastinstall"`
}

type State struct {
	On             bool      `json:"on"`
	Hue            uint16    `json:"hue,omitempty"`
	Effect         string    `json:"effect,omitempty"`
	Bri            uint8     `json:"bri,omitempty"`
	Sat            uint8     `json:"sat,omitempty"`
	CT             uint16    `json:"ct,omitempty"`
	XY             []float32 `json:"xy,omitempty"`
	Alert          string    `json:"alert,omitempty"`
	TransitionTime uint16    `json:"transitiontime,omitempty"`
	Reachable      bool      `json:"reachable,omitempty"`
	ColorMode      string    `json:"colormode,omitempty"`
	Mode           string    `json:"mode,omitempty"`
}
