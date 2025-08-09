package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	"bot/config"
	ev "bot/events"
	_ "bot/messaging/plugins"
	sql "bot/sql"
	"bot/utils"
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

	err := client.Connect()
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	utils.PairClient(ctx, client, config.AppConfig)
	utils.SetClient(client)
	utils.PortServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	client.Disconnect()
}
