package sqlite

import (
	"time"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Site struct {
	gorm.Model
	LastCrawled time.Time
	Address string
}

type Page struct {
	gorm.Model
	Site Site
	SiteID int
	Address string
}

type PageText struct {
	gorm.Model
	Page Page
	PageID int
	Text string
}