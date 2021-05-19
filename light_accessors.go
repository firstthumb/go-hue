package hue

func (l *Light) GetID() int {
	return *l.ID
}

func (l *Light) GetName() string {
	if l == nil || l.Name == nil {
		return ""
	}
	return *l.Name
}

func (l *Light) GetType() string {
	if l == nil || l.Type == nil {
		return ""
	}
	return *l.Type
}

// State

func (l *Light) IsOn() bool {
	if l == nil || l.State == nil || l.State.On == nil {
		return false
	}
	return *l.State.On
}

func (l *Light) GetEffect() string {
	if l == nil || l.State == nil || l.State.Effect == nil {
		return ""
	}
	return *l.State.Effect
}

func (l *Light) GetBri() uint8 {
	if l == nil || l.State == nil || l.State.Bri == nil {
		return 0
	}
	return *l.State.Bri
}

func (l *Light) GetSat() uint8 {
	if l == nil || l.State == nil || l.State.Sat == nil {
		return 0
	}
	return *l.State.Sat
}

func (l *Light) GetCT() uint16 {
	if l == nil || l.State == nil || l.State.CT == nil {
		return 0
	}
	return *l.State.CT
}

func (l *Light) GetXY() []float32 {
	if l == nil || l.State == nil || l.State.XY == nil {
		return []float32{}
	}
	return l.State.XY
}

func (l *Light) GetAlert() string {
	if l == nil || l.State == nil || l.State.Alert == nil {
		return ""
	}
	return *l.State.Alert
}

func (l *Light) GetTransitionTime() uint16 {
	if l == nil || l.State == nil || l.State.TransitionTime == nil {
		return 0
	}
	return *l.State.TransitionTime
}

func (l *Light) IsReachable() bool {
	if l == nil || l.State == nil || l.State.Reachable == nil {
		return false
	}
	return *l.State.Reachable
}

func (l *Light) GetColorMode() string {
	if l == nil || l.State == nil || l.State.ColorMode == nil {
		return ""
	}
	return *l.State.ColorMode
}

func (l *Light) GetMode() string {
	if l == nil || l.State == nil || l.State.Mode == nil {
		return ""
	}
	return *l.State.Mode
}
