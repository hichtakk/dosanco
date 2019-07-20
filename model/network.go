package model

import (
	"net"
	"strconv"
	"strings"
)

// IPv4Network represents network specification
type IPv4Network struct {
	Model
	CIDR           string        `json:"cidr" validate:"required" gorm:"unique;not null"`
	Description    string        `json:"description"`
	SupernetworkID uint          `json:"supernet_id" validate:"required"`
	Subnetworks    []IPv4Network `json:"subnets,omitempty"`
	Reserved       bool          `gorm:"default:false"`
	Allocations    []IPv4Allocation
}

//type IPv6Network struct{}

// Vlan
type Vlan struct {
	Model
	Description   string `json:"description"`
	IPv4NetworkID uint
	//IPv4Network   IPv4Network
	//IPv6NetworkID uint
	//IPv6Network IPv6Network
}

// IPAM
type IPv4Allocation struct {
	Model
	Address       string `json:"address" gorm:"unique;not null" validate:"required"`
	Name          string `json:"name" gorm:"not null" validate:"required"`
	Description   string `json:"description"`
	IPv4NetworkID uint   `json:"ipv4_network_id" gorm:"not null" validate:"required" sql:"type:integer REFERENCES ipv4_networks(id)"`
}

// type IPv6Allocation struct {}

// GetNetwork returns net.IPNet instance address of IPv4Network
func (n IPv4Network) GetNetwork() *net.IPNet {
	_, ipv4Net, _ := net.ParseCIDR(n.CIDR)
	return ipv4Net
}

// GetPrefixLength returns prefix size of the network
func (n IPv4Network) GetPrefixLength() int {
	slice := strings.Split(n.CIDR, "/")
	len, _ := strconv.Atoi(slice[1])
	return len
}
