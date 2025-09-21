package main

import (
	"database/sql"
	"log"

	"github.com/Ma-hiru/simplebank/api"
	db "github.com/Ma-hiru/simplebank/db/sqlc"
	"github.com/Ma-hiru/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	var config, err = util.LoadConfig(".", "app")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	{
		var conn, err = sql.Open(config.DBDriver, config.DBSource)
		if err != nil {
			log.Fatal("cannot connect to db:", err)
		}
		var store = db.NewStore(conn)
		var server = api.NewServer(store)

		err = server.Start(config.ServerAddress)
		if err != nil {
			log.Fatal("cannot start server:", err)
		}
	}

}
