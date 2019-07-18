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
}

//type IPv6Network struct{}

// Vlan
type Vlan struct {
	Model
	Description   string `json:"description"`
	IPv4NetworkID uint
	IPv4Network   IPv4Network
	//IPv6NetworkID uint
	//IPv6Network IPv6Network
}

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
