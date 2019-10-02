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
	RackID          uint            `json:"rack_id"`
	Description     string          `json:"description"`
	IPv4Allocations IPv4Allocations `json:"ipv4_allocations"`
	Rack            Rack
	//Location        string          `json:"location"`
	//MountUnit
	//Configuration   Configuration   `json:"configuration" `
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

/*
type CPU struct {
	Model string `json:"model"`
}

type Memory struct {
	Capacity int `json:"capacity"`
}

type Drive struct {
	Capacity int `json:"capacity"`
}

type NIC struct {
	Model string `json:"nic"`
}

type Accelerator struct {
	Model string `json:"accelerator"`
}

type Configuration struct {
	CPUs         []CPU         `json:"cpu"`
	Memories     []Memory      `json:"memory"`
	Storage      []Drive       `json:"storage"`
	NICs         []NIC         `json:"nic"`
	Accelerators []Accelerator `json:"accelerator"`
}
*/
