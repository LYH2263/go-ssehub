package ssehub

import "github.com/LYH2263/go-ssehub/internal/event"

func NewDefaultEncoder() event.Encoder {
	return event.DefaultEncoder{}
}
