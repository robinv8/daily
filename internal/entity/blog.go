package entity

type Blog struct {
	Id        string `xorm:"not null pk VARCHAR(20) id" json:"id"`
	BlogUrl   string `xorm:"not null VARCHAR(100) blog_url" json:"blog_url"`
	Logo      string `xorm:"VARCHAR(100) logo" json:"logo"`
	BlogName  string `xorm:"VARCHAR(100) blog_name" json:"blog_name"`
	BlogCover string `xorm:"VARCHAR(100) blog_cover" json:"blog_cover"`
}

func (Blog) TableName() string {
	return "blog"
}
