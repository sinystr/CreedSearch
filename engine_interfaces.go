package main

import "./models"

type DatabaseEngine interface {
	ShouldSiteBeCrawled(site string) bool
	GetRecordForSite(siteAddress string) models.Site
	UpdateSiteRecord(siteRecord models.Site)
}
	
type SearchEngine interface {
	search(site models.Site, text string) []models.Page
}

type CrawlingEngine interface {
	CrawlSite(site string) (models.Site, error)
}