package proper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleanParam(t *testing.T) {
	parameters := make(map[string]string)
	parameters["boolean"] = "true"
	parameters["integer"] = "1"
	parameters["bad"] = "bad"

	emptyProperties := &Properties{}
	properties := NewProperties(parameters)

	assert.Equal(t, emptyProperties.BooleanParam("boolean", false), false)
	assert.Equal(t, properties.BooleanParam("boolean", false), true)
	assert.Equal(t, properties.BooleanParam("integer", false), true)
	assert.Equal(t, properties.BooleanParam("bad", false), false)
}

func TestFloatParam(t *testing.T) {
	parameters := make(map[string]string)
	parameters["integer"] = "11"
	parameters["float"] = "1.2"
	parameters["bad"] = "bad"

	emptyProperties := &Properties{}
	properties := NewProperties(parameters)

	assert.Equal(t, emptyProperties.FloatParam("float", 0), float64(0))
	assert.Equal(t, properties.FloatParam("integer", 0), float64(11))
	assert.Equal(t, properties.FloatParam("float", 0), float64(1.2))
	assert.Equal(t, properties.FloatParam("bad", 0), float64(0))
}

func TestIntegerParam(t *testing.T) {
	parameters := make(map[string]string)
	parameters["integer"] = "11"
	parameters["float"] = "1.2"
	parameters["bad"] = "bad"

	emptyProperties := &Properties{}
	properties := NewProperties(parameters)

	assert.Equal(t, emptyProperties.IntegerParam("integer", 0), 0)
	assert.Equal(t, properties.IntegerParam("integer", 0), 11)
	assert.Equal(t, properties.IntegerParam("float", 0), 0)
	assert.Equal(t, properties.IntegerParam("bad", 0), 0)
}

func TestStringParam(t *testing.T) {
	parameters := make(map[string]string)
	parameters["string"] = "value"

	emptyProperties := &Properties{}
	properties := NewProperties(parameters)

	assert.Equal(t, emptyProperties.StringParam("string", ""), "")
	assert.Equal(t, properties.StringParam("string", ""), "value")
}
