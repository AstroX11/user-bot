package messaging

import (
	"fmt"
	"log"
	"time"

	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func Ping(msg *events.Message) {
	start := time.Now()

	id, err := utils.SendMessage(msg.Info.Chat, "🏓 Pong!")
	if err != nil {
		log.Panicln(err)
	}

	duration := time.Since(start)
	_ = utils.EditMessage(msg.Info.Chat, id, fmt.Sprintf("```Pong (%v)```", duration.Round(time.Millisecond)))
}
