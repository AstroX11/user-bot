package types

import "go.mau.fi/whatsmeow/types/events"

type Command struct {
	Name     string
	Category string
	FromMe   bool
	IsGroup  bool
	Handler  func(msg *events.Message, args []string)
}
