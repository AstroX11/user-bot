package utils

import (
	"context"
	"log"
	"time"

	"go.mau.fi/whatsmeow"

	app "github.com/AstroX11/user-bot/types"
)

func PairClient(ctx context.Context, client *whatsmeow.Client, config app.Config) {
	if client.Store.ID == nil {
		time.Sleep(2 * time.Second)

		log.Println("Pairing new session with phone number:", config.UserPN)

		code, err := client.PairPhone(ctx, config.UserPN, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			log.Fatal("Pairing failed:", err)
		}

		log.Println("Pair Code:", code)
	}
}
