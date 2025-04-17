package entity

import "time"

type SiteInfo struct {
	Id          string    `xorm:"not null pk VARCHAR(20) id" json:"id"`
	Title       string    `xorm:"not null VARCHAR(255) title" json:"title"`
	Keywords    string    `xorm:"VARCHAR(255) keywords" json:"keywords"`
	Description string    `xorm:"TEXT description" json:"description"`
	OriginUrl   string    `xorm:"not null VARCHAR(512) origin_url" json:"origin_url"`
	ImageUrl    string    `xorm:"not null VARCHAR(512) image_url" json:"image_url"`
	CreatedAt   time.Time `xorm:"not null TIMESTAMP created_at" json:"created_at"`
}

func (SiteInfo) TableName() string {
	return "site_info"
}
