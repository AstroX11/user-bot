package test

import (
	"log"

	"bot/sql"
)

func Test() {
	log.Println(sql.GetPrefix())
}
