package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  // mysql/mariadb
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite

	"github.com/hichikaw/dosanco/config"
	"github.com/hichikaw/dosanco/model"
)

var (
	db     *gorm.DB
	err    error
	schema []string
)

// Init initializes database connection and ORM
func Init(c config.Config) {
	schema = strings.Split(c.DB.URL, "://")
	if schema[0] == "sqlite" {
		db, err = gorm.Open("sqlite3", schema[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect database: "+err.Error())
			os.Exit(255)
		}
		db.Exec("PRAGMA foreign_keys = ON;")
	} else if schema[0] == "mysql" {
		db, err = gorm.Open("mysql", schema[1]+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			os.Exit(255)
		}
	}

	if c.Feature.Network {
		initNetwork()
		initIPAllocation()
	}
	if c.Feature.Host {
		initHost()
	}
	if c.Feature.DataCenter {
		initDataCenter()
	}
}

func initDataCenter() {
	db.AutoMigrate(&model.DataCenter{})
	db.AutoMigrate(&model.Floor{})
	db.AutoMigrate(&model.Hall{})
	db.AutoMigrate(&model.RackRow{})
	db.AutoMigrate(&model.Rack{})
	db.AutoMigrate(&model.UPS{})
	db.AutoMigrate(&model.PDU{})
	db.AutoMigrate(&model.RackPDU{})
}

func initNetwork() {
	db.AutoMigrate(&model.IPv4Network{})
	db.AutoMigrate(&model.Vlan{})
	var rootNetwork model.IPv4Network
	if result := db.Take(&rootNetwork, "id=1"); result.Error != nil {
		rootNetwork.ID = 1
		rootNetwork.CIDR = "0.0.0.0/0"
		rootNetwork.Description = "Root"
		rootNetwork.Reserved = true
		db.Create(&rootNetwork)
	}
	/*
		// reserve RFC6890
		var hostNetwork model.IPv4Network
		if result := db.Take(&hostNetwork, "id=2"); result.Error != nil {
			hostNetwork.ID = 2
			hostNetwork.CIDR = "0.0.0.0/8"
			hostNetwork.SupernetworkID = 1
			hostNetwork.Description = "RFC 1122,6890: This host on this network"
			hostNetwork.Reserved = true
			db.Create(&hostNetwork)
		}
		var loopbackNetwork model.IPv4Network
		if result := db.Take(&loopbackNetwork, "id=3"); result.Error != nil {
			loopbackNetwork.ID = 3
			loopbackNetwork.CIDR = "127.0.0.0/8"
			loopbackNetwork.SupernetworkID = 1
			loopbackNetwork.Description = "RFC 6890: Loopback"
			loopbackNetwork.Reserved = true
			db.Create(&loopbackNetwork)
		}
		var linkLocalNetwork model.IPv4Network
		if result := db.Take(&linkLocalNetwork, "id=4"); result.Error != nil {
			linkLocalNetwork.ID = 4
			linkLocalNetwork.CIDR = "169.254.0.0/16"
			linkLocalNetwork.SupernetworkID = 1
			linkLocalNetwork.Description = "RFC 3927: Link-local"
			linkLocalNetwork.Reserved = true
			db.Create(&linkLocalNetwork)
		}
		var protoAssignNetwork model.IPv4Network
		if result := db.Take(&protoAssignNetwork, "id=5"); result.Error != nil {
			protoAssignNetwork.ID = 5
			protoAssignNetwork.CIDR = "192.0.0.0/24"
			protoAssignNetwork.SupernetworkID = 1
			protoAssignNetwork.Description = "RFC 6890: IETF protocol assignment"
			protoAssignNetwork.Reserved = true
			db.Create(&protoAssignNetwork)
		}
		var testNetwork model.IPv4Network
		if result := db.Take(&testNetwork, "id=6"); result.Error != nil {
			testNetwork.ID = 6
			testNetwork.CIDR = "192.0.2.0/24"
			testNetwork.SupernetworkID = 1
			testNetwork.Description = "RFC 5737,6890: TEST-NET-1"
			testNetwork.Reserved = true
			db.Create(&testNetwork)
		}
		var testNetwork2 model.IPv4Network
		if result := db.Take(&testNetwork2, "id=7"); result.Error != nil {
			testNetwork2.ID = 7
			testNetwork2.CIDR = "198.51.100.0/24"
			testNetwork2.SupernetworkID = 1
			testNetwork2.Description = "RFC 5737,6890: TEST-NET-2"
			testNetwork2.Reserved = true
			db.Create(&testNetwork2)
		}
		var testNetwork3 model.IPv4Network
		if result := db.Take(&testNetwork3, "id=8"); result.Error != nil {
			testNetwork3.ID = 8
			testNetwork3.CIDR = "203.0.113.0/24"
			testNetwork3.SupernetworkID = 1
			testNetwork3.Description = "RFC 5737,6890: TEST-NET-3"
			testNetwork3.Reserved = true
			db.Create(&testNetwork3)
		}
	*/
}

func initIPAllocation() {
	if schema[0] == "sqlite" {
		// refer https://github.com/jinzhu/gorm/issues/765
		db.AutoMigrate(&model.IPv4Allocation{})
	} else {
		db.AutoMigrate(&model.IPv4Allocation{}).AddForeignKey("ipv4_network_id", "ipv4_networks(id)", "CASCADE", "CASCADE")
	}
}

func initHost() {
	db.AutoMigrate(&model.Host{})
	db.AutoMigrate(&model.HostGroup{})
}

// GetDB returns database pointer
func GetDB() *gorm.DB {
	return db
}
