package utils

import (
	"go.mau.fi/whatsmeow/proto/waE2E"
)

func ExtractText(msg *waE2E.Message) string {
	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()

	case msg.ExtendedTextMessage != nil:
		return msg.ExtendedTextMessage.GetText()

	case msg.ImageMessage != nil && msg.ImageMessage.Caption != nil:
		return msg.ImageMessage.GetCaption()

	case msg.VideoMessage != nil && msg.VideoMessage.Caption != nil:
		return msg.VideoMessage.GetCaption()

	case msg.ProtocolMessage != nil && msg.ProtocolMessage.EditedMessage != nil:
		return ExtractText(msg.ProtocolMessage.EditedMessage)
	}
	return ""
}

func GetQuotedMessage(msg *waE2E.ContextInfo) {

}
