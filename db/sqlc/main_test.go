package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Ma-hiru/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var config, err = util.LoadConfig("../..", "app")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	{
		var conn, err = sql.Open(config.DBDriver, config.DBSource)
		testDB = conn

		if err != nil {
			log.Fatal("cannot connect to db:", err)
		}
		testQueries = New(conn)
	}

	os.Exit(m.Run())
}
