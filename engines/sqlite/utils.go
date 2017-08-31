package sqlite

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// OpenDatabase open the database
// Note: if the database does not exist it creates it
func OpenDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "database.db")

	if err != nil {
	  panic("failed to connect database")
	}
  
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Page{})
	db.AutoMigrate(&PageText{})
	
	return db
}
