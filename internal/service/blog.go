package service

import (
	"daily/internal/db"
	"daily/internal/entity"
)

func Blog() []entity.Blog {
	var dailyDB = db.NewDB()

	var blog []entity.Blog
	err := dailyDB.Find(&blog)
	if err != nil {
		panic(err)
	}

	return blog
}
