package main

import (
	"database/sql"
	"log"

	"github.com/stripe/stripe-go/v81"

	"github.com/Yelsnik/e-commerce-api/api"
	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	stripe.Key = config.StripeSecretKey

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
