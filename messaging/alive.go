package messaging

import (
	"fmt"
	"time"

	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func Alive(msg *events.Message) {
	uptime := time.Since(time.Now().Add(-time.Minute * 5))
	response := fmt.Sprintf("✅ Bot is alive!\n⏰ Uptime: %v\n🤖 Status: Running", uptime.Round(time.Second))

	_, _ = utils.SendMessage(msg.Info.Chat, response)
}
