package model

import (
	"fmt"
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
	Subnetworks    IPv4Networks `json:"subnets,omitempty"`
	Reserved       bool          `json:"reserved" gorm:"default:false"`
	Allocations    []IPv4Allocation `json:"allocations,omitempty"`
}

func (n IPv4Network) Write() {
	fmt.Printf("# Network Data\n")
	fmt.Printf(" ID:             %-3d\n", n.ID)
	fmt.Printf(" CIDR:           %v\n", n.CIDR)
	fmt.Printf(" Description:    %v\n", n.Description)
	fmt.Printf(" SupernetworkID: %d\n\n", n.SupernetworkID)
	if len(n.Subnetworks) > 0 {
		fmt.Println("# Subnetworks")
		for _, s := range n.Subnetworks {
			fmt.Printf(" %-15v %v\n", s.CIDR, s.Description)
		}
	}
	if len(n.Allocations) > 0 {
		fmt.Println("# IP Allocations")
		for _, a := range n.Allocations {
			fmt.Printf(" %-15v %v, %v\n", a.Address, a.Name, a.Description)
		}
	}
}

//type IPv6Network struct{}

// Vlan
type Vlan struct {
	Model
	Description   string `json:"description"`
	IPv4NetworkID uint	`json:"ipv4_network_id" gorm:"not null"`
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

func (n IPv4Network) GetNetworkAddress() string {
	ipNet := n.GetNetwork()
	return ipNet.IP.String()
}

// GetPrefixLength returns prefix size of the network
func (n IPv4Network) GetPrefixLength() int {
	slice := strings.Split(n.CIDR, "/")
	len, _ := strconv.Atoi(slice[1])
	return len
}

func (n IPv4Network) Contains(addr string) bool {
	ip := net.ParseIP(addr)
	ipnet := n.GetNetwork()
	return ipnet.Contains(ip)
}

func (a IPv4Allocation) GetAddressInteger() uint32 {
	binAddress := ""
	octets := strings.Split(a.Address, ".")
	for _, octet := range octets {
		i, _ := strconv.Atoi(octet)
		binAddress = binAddress + fmt.Sprintf("%08b", i)
	}
	u64, _ := strconv.ParseUint(binAddress, 2, 32)

	return uint32(u64)
}


// 
type IPv4Networks []IPv4Network
func (n IPv4Networks) Len() int {
	return len(n)
}

func (n IPv4Networks) Less(i, j int) bool {
	iIP := n[i].GetNetworkAddress()
	jIP := n[j].GetNetworkAddress()
	ipCompStr := []string{iIP, jIP}
	ipComp := []uint32{}
	for _, i := range ipCompStr {
		binAddress := ""
		octets := strings.Split(i, ".")
		for _, octet := range octets {
			i, _ := strconv.Atoi(octet)
			binAddress = binAddress + fmt.Sprintf("%08b", i)
		}
		u64, _ := strconv.ParseUint(binAddress, 2, 32)
		ipComp = append(ipComp, uint32(u64))
	}
	if ipComp[0] < ipComp[1] {
		return true
	} else if ipComp[0] == ipComp[1] {
		iLen := n[i].GetPrefixLength()
		jLen := n[j].GetPrefixLength()
		return iLen < jLen
	} else {
		return false
	}
}

func (n IPv4Networks) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

//
type IPv4Allocations []IPv4Allocation

func (a IPv4Allocations) Len() int {
	return len(a)
}

func (a IPv4Allocations) Less(i, j int) bool {
	return a[i].GetAddressInteger() < a[j].GetAddressInteger()
}

func (a IPv4Allocations) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
