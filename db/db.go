package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/hichikaw/dosanco/model"
)

var (
	db *gorm.DB
	err error
)

func init() {
	db, err = gorm.Open("sqlite3", "dosanco.db")
	if err != nil {
		panic("failed to connect database")
	}
	/*
	db.AutoMigrate(&model.DataCenter{})
	db.AutoMigrate(&model.Floor{})
	db.AutoMigrate(&model.Hall{})
	db.AutoMigrate(&model.RackRow{})
	db.AutoMigrate(&model.Rack{})
	db.AutoMigrate(&model.UPS{})
	db.AutoMigrate(&model.PDU{})
	db.AutoMigrate(&model.RackPDU{})
	*/
	db.AutoMigrate(&model.IPv4Network{})
}

func GetDB() *gorm.DB {
	return db
}