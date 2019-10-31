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

// Take returns DataCenter matches specified ID
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
	Name         string `gorm:"type:varchar(16);unique_dc_floor" json:"name"`
	DataCenterID uint   `gorm:"unique_dc_floor" json:"datacenter_id"`
	DataCenter   DataCenter
	Halls        Halls `json:"halls,omitempty"`
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

// Floors represents slice of Floor
type Floors []Floor

func (f Floors) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(f, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %-5s   %-10s\n", "ID", "DC", "Name")
		for _, floor := range f {
			fmt.Printf("%3d   %5s   %-10s\n", floor.ID, floor.DataCenter.Name, floor.Name)
		}
	} else {
		fmt.Printf("%-5s   %-10s\n", "DC", "Name")
		for _, floor := range f {
			fmt.Printf("%5s   %-10s\n", floor.DataCenter.Name, floor.Name)
		}
	}
}

// Len returns number of Floors
func (f Floors) Len() int {
	return len(f)
}

// Take returns Floor specified by ID
func (f Floors) Take(id uint) (*Floor, error) {
	for _, floor := range f {
		if floor.ID == id {
			return &floor, nil
		}
	}

	return &Floor{}, fmt.Errorf("no floor found for '%v'", id)
}

// Hall represents data hall in datacenter.
type Hall struct {
	Model
	Name     string    `gorm:"type:varchar(16)" json:"name"`
	FloorID  uint      `json:"floor_id"`
	Floor    Floor     `json:"floor"`
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
		fmt.Printf(" FloorID: %v\n", h.FloorID)
	}
}

// Halls represents slice of Hall
type Halls []Hall

func (h Halls) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(h, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%3s   %7s   %10s\n", "ID", "Floor", "Name")
		for _, hall := range h {
			fmt.Printf("%3d   %7v   %10s\n", hall.ID, hall.Floor.Name, hall.Name)
		}
	}
}

// Take returns Hall specified by ID
func (h Halls) Take(id uint) (*Hall, error) {
	for _, hall := range h {
		if hall.ID == id {
			return &hall, nil
		}
	}

	return &Hall{}, fmt.Errorf("no hall found for '%v'", id)
}

// RackRow represents row of racks in data hall.
type RackRow struct {
	Model
	Name   string `gorm:"type:varchar(16)" json:"name"`
	HallID uint   `json:"hall_id"`
	Hall   Hall   `json:"hall"`
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
		fmt.Printf(" Hall: %v\n", r.Hall.Name)
	}
}

// RackRows represents slice of RackRow
type RackRows []RackRow

func (r RackRows) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %-5s   %-5s   %7s   %10s\n", "ID", "DC", "Floor", "Hall", "Name")
		for _, row := range r {
			fmt.Printf("%3d   %5s   %5s   %7s   %10s\n", row.ID, row.Hall.Floor.DataCenter.Name, row.Hall.Floor.Name, row.Hall.Name, row.Name)
		}
	} else {
		fmt.Printf("%-5s   %-5s   %7s   %10s\n", "DC", "Floor", "Hall", "Name")
		for _, row := range r {
			fmt.Printf("%5s   %5s   %7s   %10s\n", row.Hall.Floor.DataCenter.Name, row.Hall.Floor.Name, row.Hall.Name, row.Name)
		}
	}
}

// Take returns RackRow specified by ID
func (r RackRows) Take(id uint) (*RackRow, error) {
	for _, row := range r {
		if row.ID == id {
			return &row, nil
		}
	}

	return &RackRow{}, fmt.Errorf("no rack row found for '%v'", id)
}

// Rack represents each rack in row.
type Rack struct {
	Model
	Name        string   `gorm:"type:varchar(16)" json:"name"`
	RowID       uint     `json:"row_id"`
	RackPDUs    RackPDUs `json:"rack_pdus,omitempty"`
	Description string   `gorm:"type:varchar(255)" json:"description"`
	RackRow     RackRow  `json:"row"`
}

func (r Rack) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# Rack\n")
		fmt.Printf(" ID:    %d\n", r.ID)
		fmt.Printf(" Name:  %v\n", r.Name)
		fmt.Printf(" Row:   %v\n", r.RackRow.Name)
	}
}

// GetLocationPath returns location of rack including datacenter, floor, hall and row
func (r Rack) GetLocationPath() string {
	row := r.RackRow.Name
	hall := r.RackRow.Hall.Name
	floor := r.RackRow.Hall.Floor.Name
	dc := r.RackRow.Hall.Floor.DataCenter.Name

	return fmt.Sprintf("%v/%v/%v/%v/%v", dc, floor, hall, row, r.Name)
}

// Racks represents slice of Rack
type Racks []Rack

func (r Racks) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%-3s   %-5s   %-5s   %-5s   %-5s   %10s   %s\n", "ID", "DC", "Floor", "Hall", "Row", "Name", "Description")
		for _, rack := range r {
			fmt.Printf("%3d   %5v   %5v   %5s   %5s   %10s   %s\n", rack.ID, rack.RackRow.Hall.Floor.DataCenter.Name, rack.RackRow.Hall.Floor.Name, rack.RackRow.Hall.Name, rack.RackRow.Name, rack.Name, rack.Description)
		}
	} else {
		fmt.Printf("%-5s   %-5s   %-5s   %-5s   %10s   %s\n", "DC", "Floor", "Hall", "Row", "Name", "Description")
		for _, rack := range r {
			fmt.Printf("%5v   %5v   %5s   %5s   %10s   %s\n", rack.RackRow.Hall.Floor.DataCenter.Name, rack.RackRow.Hall.Floor.Name, rack.RackRow.Hall.Name, rack.RackRow.Name, rack.Name, rack.Description)
		}
	}
}

// UPS represents redundant power source
type UPS struct {
	Model
	Name         string     `gorm:"type:varchar(16)" json:"name"`
	DataCenterID uint       `json:"datacenter_id"`
	Description  string     `gorm:"type:varchar(255)" json:"description"`
	DataCenter   DataCenter `json:"datacenter"`
}

func (u UPS) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(u, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# UPS\n")
		fmt.Printf(" ID:     %d\n", u.ID)
		fmt.Printf(" Name:   %v\n", u.Name)
		fmt.Printf(" DataCenter: %v\n", u.DataCenter.Name)
	}
}

// UPSs represents slice of UPS
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

// Take returns UPS specified by ID
func (u UPSs) Take(id uint) (*UPS, error) {
	for _, ups := range u {
		if ups.ID == id {
			return &ups, nil
		}
	}

	return &UPS{}, fmt.Errorf("no ups found for '%v'", id)
}

// RowPDU represents power distribution unit on data hall
type RowPDU struct {
	Model
	Name           string `gorm:"type:varchar(16)" json:"name"`
	PrimaryUPSID   uint   `json:"primary_ups_id,omitempty"`
	SecondaryUPSID uint   `json:"secondary_ups_id,omitempty"`
	Description    string `gorm:"type:varchar(255)" json:"description"`
	PrimaryUPS     UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS   UPS    `json:"secondary_ups,omitempty"`
}

func (p RowPDU) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# DataCenter RowPDU\n")
		fmt.Printf(" ID:            %d\n", p.ID)
		fmt.Printf(" Name:          %v\n", p.Name)
		fmt.Printf(" Primary UPS:   %v\n", p.PrimaryUPS.Name)
		fmt.Printf(" Secondary UPS: %v\n", p.SecondaryUPS.Name)
	}
}

// RowPDUs represents slice of RowPDU
type RowPDUs []RowPDU

func (p RowPDUs) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%3s   %-12s   %-16s   %-16s   %-32s\n", "ID", "Name", "InputUPS#1", "InputUPS#2", "Description")
		for _, pdu := range p {
			fmt.Printf("%3d   %-12s   %-16v   %-16v   %-32v\n", pdu.ID, pdu.Name, pdu.PrimaryUPS.Name, pdu.SecondaryUPS.Name, pdu.Description)
		}
	} else {
		fmt.Printf("%-12s   %-16s   %-16s   %-32v\n", "Name", "InputUPS#1", "InputUPS#2", "Description")
		for _, pdu := range p {
			fmt.Printf("%-12s   %-16v   %-16v   %-32v\n", pdu.Name, pdu.PrimaryUPS.Name, pdu.SecondaryUPS.Name, pdu.Description)
		}
	}
}

// Take returns RowPDU specified by ID
func (p RowPDUs) Take(id uint) (*RowPDU, error) {
	for _, pdu := range p {
		if pdu.ID == id {
			return &pdu, nil
		}
	}

	return &RowPDU{}, fmt.Errorf("no pdu found for '%v'", id)
}

// RackPDU represents power distribution unit installed inside of rack
type RackPDU struct {
	Model
	Name           string  `gorm:"type:varchar(64)" json:"name"`
	Description    string  `gorm:"type:varchar(255)" json:"description"`
	PrimaryPDUID   uint    `gorm:"column:primary_pdu_id" json:"primary_pdu_id,omitempty"`
	SecondaryPDUID uint    `gorm:"column:secondary_pdu_id" json:"secondary_pdu_id,omitempty"`
	PrimaryPDU     RowPDU  `json:"primary_pdu"`
	SecondaryPDU   *RowPDU `json:"secondary_pdu"`
	Host           *Host   `json:"host,omitempty"`
}

func (p RackPDU) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(p, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("# Rack PDU\n")
		fmt.Printf(" ID:          %d\n", p.ID)
		fmt.Printf(" Name:        %v\n", p.Name)
		fmt.Printf(" Input#1:     %v\n", p.PrimaryPDU.Name)
		if p.SecondaryPDU != nil {
			fmt.Printf(" Input#2:     %v\n", p.SecondaryPDU.Name)
		} else {
			fmt.Printf(" Input#2:     %v\n", "-")
		}
		fmt.Printf(" Description: %v\n", p.Description)
	}
}

// RackPDUs represents slice of RackPDU
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
		fmt.Printf("%-32s   %-12s   %-12s   %-10s\n", "Name", "Input#1", "Input#2", "Rack")
		for _, pdu := range p {
			location := "-"
			if pdu.Host.ID != 0 {
				location = pdu.Host.Rack.GetLocationPath()
			}
			secondRowPDU := "-"
			if pdu.SecondaryPDU != nil {
				secondRowPDU = pdu.SecondaryPDU.Name
			}
			fmt.Printf("%32s   %12s   %12s   %10s\n", pdu.Name, pdu.PrimaryPDU.Name, secondRowPDU, location)
		}
	}
}
