package main

import (
	"log"

	"github.com/crema-labs/gitgo/pkg/server"
	"github.com/crema-labs/gitgo/pkg/store"
)

func main() {
	db, err := store.NewSQLiteStore("./grants.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := server.NewServer(db)
	log.Fatal(s.Start(":8080"))
}
