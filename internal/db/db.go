package db

import (
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

func NewDB() *xorm.Engine {
	db, err := xorm.NewEngine("postgres", "user=postgres password=hellodaily dbname=daily sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
