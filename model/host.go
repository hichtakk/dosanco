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
	Group           *HostGroup      `json:"group,omitempty"`
	IPv4Allocations IPv4Allocations `json:"ipv4_allocations"`
	RackID          uint            `json:"rack_id"`
	Rack            Rack            `json:"rack,omitempty"`
	Type            string          `json:"type" gorm:"default:'generic';not null"`
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
		fmt.Printf(" Group:          %v\n", h.Group.Name)
		fmt.Printf(" Type:           %v\n", h.Type)
		fmt.Printf(" RackLocation:   %v\n", h.Rack.GetLocationPath())
		fmt.Printf(" Description:    %v\n\n", h.Description)
		fmt.Println("# IP Allocations")
		if h.IPv4Allocations.Len() > 0 {
			for _, a := range h.IPv4Allocations {
				fmt.Printf(" %-15v %-15v %v\n", a.IPv4Network.CIDR, a.Address, a.Description)
			}
		} else {
			fmt.Println(" allocation not found")
		}
	}
}

// Hosts represents list of Host.
type Hosts []Host

func (h Hosts) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(h, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%-30v   %-30v\n", "Name", "Description")
		for _, host := range h {
			fmt.Printf("%-30v   %-30v\n", host.Name, host.Description)
		}
	}
}

// HostGroup represents group for host
type HostGroup struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"deleted_at;unique_index:unique_group" json:"deleted_at,omitempty"`
	Name        string     `json:"name" validate:"required" gorm:"unique_index:unique_group;not null"`
	Description string     `json:"description"`
	Hosts       *Hosts     `json:"hosts,omitempty"`
}

// HostGroups represents slice of HostGroup
type HostGroups []HostGroup

func (g HostGroups) Write(output string) {
	if output == "json" {

	} else {
		fmt.Printf("%-15v   %-15v\n", "Name", "Description")
		for _, grp := range g {
			fmt.Printf("%15v   %15v\n", grp.Name, grp.Description)
		}
	}

}
