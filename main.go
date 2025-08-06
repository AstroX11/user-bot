package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	sql "github.com/AstroX11/user-bot/sql"

	ev "github.com/AstroX11/user-bot/events"
)

var client *whatsmeow.Client

func SetClient(c *whatsmeow.Client) {
	client = c
}

func main() {
	ctx := context.Background()

	store := sqlstore.NewWithDB(sql.Conn, "sqlite", waLog.Stdout("DB", "ERROR", true))
	defer store.Close()

	if err := store.Upgrade(ctx); err != nil {
		panic("❌ Store upgrade failed: " + err.Error())
	}

	device, _ := store.GetFirstDevice(ctx)
	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "INFO", true))
	client.AddEventHandler(ev.EventHandler)

	SetClient(client)
	client.Connect()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	client.Disconnect()
}
