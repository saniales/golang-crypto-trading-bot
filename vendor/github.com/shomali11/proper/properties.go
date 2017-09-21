package proper

import "strconv"

// NewProperties creates a new Properties object
func NewProperties(m map[string]string) *Properties {
	return &Properties{propertyMap: m}
}

// Properties is a string map decorator
type Properties struct {
	propertyMap map[string]string
}

// StringParam attempts to look up a string value by key. If not found, return the default string value
func (p *Properties) StringParam(key string, defaultValue string) string {
	value, ok := p.propertyMap[key]
	if !ok {
		return defaultValue
	}
	return value
}

// BooleanParam attempts to look up a boolean value by key. If not found, return the default boolean value
func (p *Properties) BooleanParam(key string, defaultValue bool) bool {
	value, ok := p.propertyMap[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return integerValue
}

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
func (p *Properties) IntegerParam(key string, defaultValue int) int {
	value, ok := p.propertyMap[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return integerValue
}

// FloatParam attempts to look up a float value by key. If not found, return the default float value
func (p *Properties) FloatParam(key string, defaultValue float64) float64 {
	value, ok := p.propertyMap[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return integerValue
}
