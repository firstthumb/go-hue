package hue

func (l *Light) GetName() string {
	if l == nil || l.Name == nil {
		return ""
	}
	return *l.Name
}
