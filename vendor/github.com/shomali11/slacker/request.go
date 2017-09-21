package slacker

import (
	"context"

	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
)

const (
	empty = ""
)

// NewRequest creates a new Request structure
func NewRequest(ctx context.Context, event *slack.MessageEvent, properties *proper.Properties) *Request {
	return &Request{Context: ctx, Event: event, properties: properties}
}

// Request contains the Event received and parameters
type Request struct {
	Context    context.Context
	Event      *slack.MessageEvent
	properties *proper.Properties
}

// Param attempts to look up a string value by key. If not found, return the an empty string
func (r *Request) Param(key string) string {
	return r.StringParam(key, empty)
}

// StringParam attempts to look up a string value by key. If not found, return the default string value
func (r *Request) StringParam(key string, defaultValue string) string {
	return r.properties.StringParam(key, defaultValue)
}

// BooleanParam attempts to look up a boolean value by key. If not found, return the default boolean value
func (r *Request) BooleanParam(key string, defaultValue bool) bool {
	return r.properties.BooleanParam(key, defaultValue)
}

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
func (r *Request) IntegerParam(key string, defaultValue int) int {
	return r.properties.IntegerParam(key, defaultValue)
}

// FloatParam attempts to look up a float value by key. If not found, return the default float value
func (r *Request) FloatParam(key string, defaultValue float64) float64 {
	return r.properties.FloatParam(key, defaultValue)
}
