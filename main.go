package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	ev "github.com/AstroX11/user-bot/events"
	sql "github.com/AstroX11/user-bot/sql"
)

var client *whatsmeow.Client

func main() {
	ctx := context.Background()

	store := sqlstore.NewWithDB(sql.Conn, "sqlite", waLog.Stdout("DB", "ERROR", true))
	defer store.Close()

	if err := store.Upgrade(ctx); err != nil {
		log.Fatal("Store upgrade failed:", err)
	}

	device, _ := store.GetFirstDevice(ctx)
	client = whatsmeow.NewClient(device, waLog.Stdout("Client", "INFO", true))
	client.AddEventHandler(ev.EventHandler)

	if client.Store.ID == nil {
		log.Println("Connecting to WhatsApp for pairing...")

		qrChan, _ := client.GetQRChannel(ctx)

		err := client.Connect()
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}

		select {
		case <-qrChan:
			log.Println("QR code event received — safe to pair")
		case <-time.After(3 * time.Second):
			log.Println("Timeout waiting for QR code, proceeding anyway")
		}

		log.Println("Pairing new session with phone number:", AppConfig.UserPN)

		code, err := client.PairPhone(ctx, AppConfig.UserPN, true, 1, AppConfig.UserName)
		if err != nil {
			log.Fatal("Pairing failed:", err)
		}

		log.Println("Linking code:", code)
	} else {
		err := client.Connect()
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	client.Disconnect()
}
