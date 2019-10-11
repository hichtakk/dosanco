package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Host represents server or network node.
type Host struct {
	ID              uint            `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time      `gorm:"deleted_at;unique_index:unique_hostname" json:"deleted_at,omitempty"`
	Name            string          `json:"name" validate:"required" gorm:"unique_index:unique_hostname;not null"`
	Description     string          `json:"description"`
	GroupID         uint            `json:"group_id"`
	IPv4Allocations IPv4Allocations `json:"ipv4_allocations"`
	RackID          uint            `json:"rack_id"`
	Rack            Rack            `json:"rack,omitempty"`
}

// Write does output to standard output.
func (h Host) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(h, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
	} else {
		fmt.Printf("# Host Data\n")
		fmt.Printf(" ID:             %d\n", h.ID)
		fmt.Printf(" Name:           %v\n", h.Name)
		fmt.Printf(" RackLocation:   %v\n", h.Rack.GetLocationPath())
		fmt.Printf(" Description:    %v\n\n", h.Description)
		if h.IPv4Allocations.Len() > 0 {
			fmt.Println("# IP Allocations")
			for _, a := range h.IPv4Allocations {
				fmt.Printf(" %-15v %v\n", a.Address, a.Description)
			}
		}
	}
}

// Hosts represents list of Host.
type Hosts []Host

func (h Hosts) Write() {
}

type HostGroup struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"deleted_at;unique_index:unique_group" json:"deleted_at,omitempty"`
	Name        string     `json:"name" validate:"required" gorm:"unique_index:unique_group;not null"`
	Description string     `json:"description"`
	Hosts       Hosts      `json:"omitempty"`
}
