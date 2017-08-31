package sqlite

import (
	"time"
	"../../models"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Engine defines sqlite based database engine for Creed
type Engine struct {
	dataExpirationDuration time.Duration
}

// ShouldSiteBeCrawled return wheather this site needs to be crawled
// by the the crawling engine or the data is the database is fresh enough
// to be used for indexing
func(e *Engine) ShouldSiteBeCrawled(siteAddress string) bool {
	db := OpenDatabase()
	defer db.Close()
	
	var site Site
	if err := db.Where(&Site{Address:siteAddress}).First(&site).Error; err != nil {
		return true;
	}

	return time.Now().Sub(site.LastCrawled).Hours() > e.dataExpirationDuration.Hours()
}

func(e *Engine) GetRecordForSite(siteName string) models.Site {
	return models.Site{}
}

func(e *Engine) UpdateSiteRecord(siteRecord models.Site) {
	
}

func(e *Engine) SetDataExpirationTime(dataExpirationDuration time.Duration) {
	e.dataExpirationDuration = dataExpirationDuration
}

func DefaultEngine() *Engine {
	return &Engine{dataExpirationDuration: time.Hour * 24}
}
