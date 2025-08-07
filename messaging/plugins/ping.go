package plugins

import (
	"fmt"
	"time"

	"github.com/AstroX11/user-bot/messaging"
	"github.com/AstroX11/user-bot/types"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "ping",
		Category: "System",
		FromMe:   false,
		IsGroup:  false,
		Handler:  Ping,
	})
}

func Ping(msg *events.Message, _ []string) {
	start := time.Now()
	id, _ := utils.SendMessage(msg.Info.Chat, "🏓 Pong!")
	duration := time.Since(start)
	_ = utils.EditMessage(msg.Info.Chat, id, fmt.Sprintf("```Pong (%v)```", duration.Round(time.Millisecond)))
}
