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
	Name    string `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string `gorm:"type:varchar(255)" json:"address"`
	Floors  Floors `json:"floors,omitempty"`
}

// Write outputs DC to standard output.
func (d DataCenter) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(d, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# DataCenter\n")
		fmt.Printf(" ID:          %d\n", d.ID)
		fmt.Printf(" Name:        %v\n", d.Name)
		fmt.Printf(" Address:    %v\n", d.Address)
		if d.Floors.Len() > 0 {
			fmt.Println("\n# Floors")
			for _, f := range d.Floors {
				fmt.Printf(" %v\n", f.Name)
			}
		}
	}
}

// Floor represents datacenter floor or area.
type Floor struct {
	Model
	Name         string `gorm:"type:varchar(16);unique_index" json:"name"`
	DataCenterID uint   `json:"datacenter_id"`
	Halls        Halls  `json:"halls,omitempty"`
}

func (f Floor) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(f, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# Floor\n")
		fmt.Printf(" ID:           %d\n", f.ID)
		fmt.Printf(" Name:         %v\n", f.Name)
		fmt.Printf(" DataCenterID: %v\n", f.DataCenterID)
		/*
			if f.Halls.Len() > 0 {
			}
		*/
	}
}

type Floors []Floor

func (f Floors) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(f, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %3s   %-10s\n", "ID", "DC", "Name")
		for _, floor := range f {
			fmt.Printf("%3d   %3d   %-10s\n", floor.ID, floor.DataCenterID, floor.Name)
		}
	} else {
		fmt.Printf("%3s   %-10s\n", "DC", "Name")
		for _, floor := range f {
			fmt.Printf("%3d   %-10s\n", floor.DataCenterID, floor.Name)
		}
	}
}

func (f Floors) Len() int {
	return len(f)
}

// Hall represents data hall in datacenter.
type Hall struct {
	Model
	Name     string    `gorm:"type:varchar(16)" json:"name"`
	Type     string    `gorm:"type:varchar(10)" json:"type"`
	FloorID  uint      `json:"floor_id"`
	RackRows []RackRow `json:"rows,omitempty"`
}

func (h Hall) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(h, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# Hall\n")
		fmt.Printf(" ID:      %d\n", h.ID)
		fmt.Printf(" Name:    %v\n", h.Name)
		fmt.Printf(" Type:    %v\n", h.Type)
		fmt.Printf(" FloorID: %v\n", h.FloorID)
	}
}

type Halls []Hall

func (h Halls) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(h, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%3s   %7s   %10s   %6s\n", "ID", "FloorID", "Name", "Type")
		for _, hall := range h {
			fmt.Printf("%3d   %7d   %10s   %6s\n", hall.ID, hall.FloorID, hall.Name, hall.Type)
		}
	}
}

// RackRow represents row of racks in data hall.
type RackRow struct {
	Model
	Name   string `gorm:"type:varchar(16)" json:"name"`
	HallID uint   `json:"hall_id"`
	Racks  Racks  `json:"racks,omitempty"`
}

func (r RackRow) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# Rack Row\n")
		fmt.Printf(" ID:     %d\n", r.ID)
		fmt.Printf(" Name:   %v\n", r.Name)
		fmt.Printf(" HallID: %v\n", r.HallID)
	}
}

type RackRows []RackRow

func (r RackRows) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%3s   %7s   %10s\n", "ID", "HallID", "Name")
		for _, row := range r {
			fmt.Printf("%3d   %7d   %10s\n", row.ID, row.HallID, row.Name)
		}
	}
}

// Rack represents each rack in row.
type Rack struct {
	Model
	Name        string   `gorm:"type:varchar(16)" json:"name"`
	RackPDUs    RackPDUs `json:"rack_pdus,omitempty"`
	Description string   `gorm:"type:varchar(255)" json:"description"`
}
type Racks []Rack

// UPS represents redundant power source
type UPS struct {
	Model
	Name        string `gorm:"type:varchar(16)" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

// PDU represents power distribution unit on data hall
type PDU struct {
	Model
	Name         string `gorm:"type:varchar(16)" json:"name"`
	PrimaryUPS   UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS UPS    `json:"secondary_ups,omitempty"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
}

// RackPDU represents power distribution unit installed inside of rack
type RackPDU struct {
	Model
	Name        string `gorm:"type:varchar(16)" json:"name"`
	Address     string `gorm:"type:varchar(15)" json:"address"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	SourcePDUs  []PDU  `json:"source_pdus,omitempty"`
}
type RackPDUs []RackPDU
