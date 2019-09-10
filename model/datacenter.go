package model

import (
	"encoding/json"
	"fmt"
)

// DataCenters represents list of DataCenter
type DataCenters []DataCenter

func (d DataCenters) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(d, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%2s	%-10s	%s\n", "ID", "Name", "Address")
		for _, dc := range d {
			fmt.Printf("%2d	%-10s	%s\n", dc.ID, dc.Name, dc.Address)
		}
	} else {
		fmt.Printf("%-10s	%s\n", "Name", "Address")
		for _, dc := range d {
			fmt.Printf("%-10s	%s\n", dc.Name, dc.Address)
		}
	}
}

// DataCenter represents datacenter building data.
type DataCenter struct {
	Model
	Name    string  `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string  `gorm:"type:varchar(255)" json:"address"`
	Floors  []Floor `json:"floors,omitempty"`
}

// Write outputs DC to standard output.
func (d DataCenter) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(d, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
	} else {
		fmt.Printf("# DataCenter\n")
		fmt.Printf(" ID:          %d\n", d.ID)
		fmt.Printf(" Name:        %v\n", d.Name)
		fmt.Printf(" Address:    %v\n", d.Address)
		/*
			if d.Floors.Len() > 0 {
				fmt.Println("# IP Allocations")
				for _, a := range h.IPv4Allocations {
					fmt.Printf(" %-15v %v\n", a.Address, a.Description)
				}
			}
		*/
	}
}

// Floor represents datacenter floor or area.
type Floor struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Halls []Hall `json:"halls,omitempty"`
}

// Hall represents data hall in datacenter.
type Hall struct {
	Model
	Name     string    `gorm:"type:varchar(16);unique_index" json:"name"`
	Type     string    `gorm:"type:varchar(10)" json:"type"`
	RackRows []RackRow `json:"rows,omitempty"`
}

// RackRow represents row of racks in data hall.
type RackRow struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Racks []Rack `json:"racks,omitempty"`
}

// Rack represents each rack in row.
type Rack struct {
	Model
	Name        string    `gorm:"type:varchar(16);unique_index" json:"name"`
	RackPDUs    []RackPDU `json:"rack_pdus,omitempty"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}

// UPS represents redundant power source
type UPS struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

// PDU represents power distribution unit on data hall
type PDU struct {
	Model
	Name         string `gorm:"type:varchar(16);unique_index" json:"name"`
	PrimaryUPS   UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS UPS    `json:"secondary_ups,omitempty"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
}

// RackPDU represents power distribution unit installed inside of rack
type RackPDU struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Address     string `gorm:"type:varchar(15)" json:"address"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	SourcePDUs  []PDU  `json:"source_pdus,omitempty"`
}
