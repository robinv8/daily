package db

import (
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

func NewDB() *xorm.Engine {
	dbHost := "101.42.239.82"
	db, err := xorm.NewEngine("postgres", "user=postgres password=xondo9-teKcob-harqez dbname=daily sslmode=disable host="+dbHost+" port=15432")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
