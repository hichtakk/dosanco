package model

import "github.com/jinzhu/gorm"

//DataCenter
type DataCenter struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string  `gorm:"type:varchar(255)" json:"address"`
	Floors  []Floor `json:"floors"`
}

type Floor struct {
	gorm.Model
	Name    string
	Address string
	Halls   []Hall
}

type Hall struct {
	gorm.Model
	Name    string
	Type    string // Data or Network
	RackRow []RackRow
}

type RackRow struct {
	gorm.Model
	Name  string
	Racks []Rack
}

type Rack struct {
	gorm.Model
	Name     string
	RackPDUs []RackPDU
}

type UPS struct {
	gorm.Model
	Name string
}

type PDU struct {
	gorm.Model
	Name         string
	Address      string
	PrimaryUPS   UPS
	SecondaryUPS UPS
}

type RackPDU struct {
	gorm.Model
	Name       string
	SourcePDUs []PDU
}
