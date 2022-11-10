package service

import (
	"daily/internal/db"
	"daily/internal/entity"
	"sort"
)

func Home() []entity.SiteInfo {
	var dailyDB = db.NewDB()

	var siteInfo []entity.SiteInfo
	err := dailyDB.Find(&siteInfo)
	if err != nil {
		panic(err)
	}
	// sort by created_at
	sort.Slice(siteInfo, func(i, j int) bool {
		return siteInfo[i].CreatedAt.Before(siteInfo[j].CreatedAt)
	})

	return siteInfo
}
