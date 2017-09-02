package sqlite

import (
	"../../models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// Engine defines sqlite based database engine for Creed
type Engine struct {
	dataExpirationDuration time.Duration
}

// ShouldSiteBeCrawled return wheather this site needs to be crawled
// by the the crawling engine or the data is the database is fresh enough
// to be used for indexing
func (e *Engine) ShouldSiteBeCrawled(siteAddress string) bool {
	db := OpenDatabase()
	defer db.Close()

	var site Site
	if err := db.Where(&Site{Address: siteAddress}).First(&site).Error; err != nil {
		return true
	}

	return time.Now().Sub(site.LastCrawled).Hours() > e.dataExpirationDuration.Hours()
}

// GetRecordForSite return site record from the database for siteName
func (e *Engine) GetRecordForSite(siteName string) models.Site {
	return models.Site{}
}

// SaveSiteRecord saves new or updates existing crawled site record to the database
func (e *Engine) SaveSiteRecord(siteRecord models.Site) {
	db := OpenDatabase()
	defer db.Close()

	user := Site{Address: siteRecord.Address}

	if db.NewRecord(user) {
		user.LastCrawled = time.Now()
		db.Create(&user)
	} else {
		db.Model(&user).Update("last_crawled", time.Now())
	}

}

// SetDataExpirationTime sets crawed site expiration time
func (e *Engine) SetDataExpirationTime(dataExpirationDuration time.Duration) {
	e.dataExpirationDuration = dataExpirationDuration
}

// DefaultEngine returns the sqlite engine using the default configuration
func DefaultEngine() *Engine {
	return &Engine{dataExpirationDuration: time.Hour * 24}
}
