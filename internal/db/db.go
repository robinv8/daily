package db

import (
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

func NewDB() *xorm.Engine {
	// dbHost := "localhost"
	dbHost := "daily-db"
	db, err := xorm.NewEngine("postgres", "user=postgres password=hellodaily dbname=daily sslmode=disable host="+dbHost+" port=5432")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
