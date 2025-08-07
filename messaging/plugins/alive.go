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
		Name:     "alive",
		Category: "System",
		FromMe:   false,
		IsGroup:  false,
		Handler:  Alive,
	})
}

func Alive(msg *events.Message, _ []string) {
	uptime := time.Since(time.Now().Add(-time.Minute * 5))
	response := fmt.Sprintf("✅ Bot is alive!\n⏰ Uptime: %v\n🤖 Status: Running", uptime.Round(time.Second))
	_, _ = utils.SendMessage(msg.Info.Chat, response)
}
