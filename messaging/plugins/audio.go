package plugins

import (
	"strings"

	"github.com/AstroX11/user-bot/messaging"
	"github.com/AstroX11/user-bot/types"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "audio",
		Category: "Test",
		FromMe:   true,
		IsGroup:  false,
		Handler:  SendTestAudio,
	})
}

func SendTestAudio(msg *events.Message, args []string) {
	audioPath := "./resources/audio.mp3"
	
	isVoice := false
	if len(args) > 0 && (strings.EqualFold(args[0], "vn") || strings.EqualFold(args[0], "ptt")) {
		isVoice = true
	}

	 utils.SendAudio(msg.Info.Chat, audioPath, isVoice)

}
