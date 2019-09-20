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

func (d DataCenters) Take(id uint) (*DataCenter, error) {
	for _, dc := range d {
		if dc.ID == id {
			return &dc, nil
		}
	}

	return &DataCenter{}, fmt.Errorf("no data center found for '%v'", id)
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
	RowID       uint     `json:"row_id"`
	RackPDUs    RackPDUs `json:"rack_pdus,omitempty"`
	Description string   `gorm:"type:varchar(255)" json:"description"`
}
type Racks []Rack

func (r Racks) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%3s   %5s   %10s   %s\n", "ID", "RowID", "Name", "Description")
		for _, rack := range r {
			fmt.Printf("%3d   %5v   %10s   %s\n", rack.ID, rack.RowID, rack.Name, rack.Description)
		}
	}
}

// UPS represents redundant power source
type UPS struct {
	Model
	Name         string `gorm:"type:varchar(16)" json:"name"`
	DataCenterID uint   `json:"datacenter_id"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
	DataCenter   DataCenter
}

type UPSs []UPS

func (u UPSs) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(u, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%-3s   %-10s   %-10s   %-s\n", "ID", "Name", "DataCenter", "Description")
		for _, ups := range u {
			fmt.Printf("%3v   %-10v   %-10s   %-s\n", ups.ID, ups.Name, ups.DataCenter.Name, ups.Description)
		}
	} else {
		fmt.Printf("%-10s   %-10s   %-s\n", "Name", "DataCenter", "Description")
		for _, ups := range u {
			fmt.Printf("%-10v   %-10s   %-s\n", ups.Name, ups.DataCenter.Name, ups.Description)
		}
	}
}

func (u UPSs) Take(id uint) (*UPS, error) {
	for _, ups := range u {
		if ups.ID == id {
			return &ups, nil
		}
	}

	return &UPS{}, fmt.Errorf("no ups found for '%v'", id)
}

// PDU represents power distribution unit on data hall
type PDU struct {
	Model
	Name           string `gorm:"type:varchar(16)" json:"name"`
	PrimaryUPSID   uint   `json:"primary_ups_id,omitempty"`
	SecondaryUPSID uint   `json:"secondary_ups_id,omitempty"`
	Description    string `gorm:"type:varchar(255)" json:"description"`
	PrimaryUPS     UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS   UPS    `json:"secondary_ups,omitempty"`
}

func (p PDU) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# DataCenter PDU\n")
		fmt.Printf(" ID:               %d\n", p.ID)
		fmt.Printf(" Name:             %v\n", p.Name)
		fmt.Printf(" Primary UPS ID:   %v\n", p.PrimaryUPSID)
		fmt.Printf(" Secondary UPS ID: %v\n", p.SecondaryUPSID)
	}
}

type PDUs []PDU

func (p PDUs) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %-12s   %-10s   %-10s\n", "ID", "Name", "InputUPS#1", "InputUPS#2")
		for _, pdu := range p {
			fmt.Printf("%3d   %-12s   %-10v   %-10v\n", pdu.ID, pdu.Name, pdu.PrimaryUPS.Name, pdu.SecondaryUPS.Name)
		}
	} else {
		fmt.Printf("%-12s   %-10s   %-10s\n", "Name", "InputUPS#1", "InputUPS#2")
		for _, pdu := range p {
			fmt.Printf("%-12s   %-10v   %-10v\n", pdu.Name, pdu.PrimaryUPS.Name, pdu.SecondaryUPS.Name)
		}
	}
}

func (p PDUs) Take(id uint) (*PDU, error) {
	for _, pdu := range p {
		if pdu.ID == id {
			return &pdu, nil
		}
	}

	return &PDU{}, fmt.Errorf("no pdu found for '%v'", id)
}

// RackPDU represents power distribution unit installed inside of rack
type RackPDU struct {
	Model
	Name           string `gorm:"type:varchar(16)" json:"name"`
	Description    string `gorm:"type:varchar(255)" json:"description"`
	PrimaryPDUID   uint   `gorm:"column:primary_pdu_id" json:"primary_pdu_id,omitempty"`
	SecondaryPDUID uint   `gorm:"column:secondary_pdu_id" json:"secondary_pdu_id,omitempty"`
	PrimaryPDU     PDU
	SecondaryPDU   PDU
}

type RackPDUs []RackPDU

func (p RackPDUs) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %-32s   %-12s   %-12s\n", "ID", "Name", "Input#1", "Input#2")
		for _, pdu := range p {
			fmt.Printf("%3d   %32s   %12s   %12s\n", pdu.ID, pdu.Name, pdu.PrimaryPDU.Name, pdu.SecondaryPDU.Name)
		}
	} else {
		fmt.Printf("%-32s   %-12s   %-12s\n", "Name", "Input#1", "Input#2")
		for _, pdu := range p {
			fmt.Printf("%32s   %12s   %12s\n", pdu.Name, pdu.PrimaryPDU.Name, pdu.SecondaryPDU.Name)
		}
	}
}
