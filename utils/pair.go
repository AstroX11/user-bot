package utils

import (
	"context"
	"log"
	"time"

	"go.mau.fi/whatsmeow"

	app "bot/types"
)

func PairClient(ctx context.Context, client *whatsmeow.Client, config app.Config) {
	if client.Store.ID == nil {
		time.Sleep(2 * time.Second)

		log.Println("Pairing new session with phone number:", config.UserPN)

		if len(config.UserPN) < 11 {
			log.Fatal("Phone number too short. Please provide a valid international phone number (e.g., 12345678900).")
		}

		code, err := client.PairPhone(ctx, config.UserPN, true, whatsmeow.PairClientChrome, "Chrome (Windows)")
		if err != nil {
			log.Fatal("Pairing failed:", err)
		}

		log.Println("Pair Code:", code)
	}
}
