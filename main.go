package main

import (
	"database/sql"
	"log"

	"github.com/iamzay/simplebank/api"
	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannnot start the server:", err)
	}
}
