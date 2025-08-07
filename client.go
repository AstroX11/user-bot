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

	ev "github.com/AstroX11/user-bot/events"
	sql "github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
)

var waClient *whatsmeow.Client

func main() {
	ctx := context.Background()

	store := sqlstore.NewWithDB(sql.Conn, "sqlite", waLog.Stdout("DB", "ERROR", true))
	defer store.Close()

	if err := store.Upgrade(ctx); err != nil {
		log.Fatal("Store upgrade failed:", err)
	}

	device, _ := store.GetFirstDevice(ctx)
	waClient = whatsmeow.NewClient(device, waLog.Stdout("Client", "INFO", true)) // fixed assignment
	waClient.AddEventHandler(ev.EventHandler)

	err := waClient.Connect()
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	utils.PairClient(ctx, waClient, AppConfig)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	waClient.Disconnect()
}
