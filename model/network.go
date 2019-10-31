package model

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPv4Network represents network specification.
type IPv4Network struct {
	Model
	CIDR           string           `json:"cidr" validate:"required" gorm:"unique;not null"`
	Description    string           `json:"description"`
	SupernetworkID uint             `json:"supernet_id"`
	Subnetworks    IPv4Networks     `json:"subnets,omitempty"`
	Reserved       bool             `json:"reserved" gorm:"default:false"`
	Allocations    []IPv4Allocation `json:"allocations,omitempty"`
}

func (n IPv4Network) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(n, "", "    ")
		fmt.Println(string(jsonBytes))
	} else {
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
			length := n.GetPrefixLength()
			addrSpace := 1 << uint(32-length)
			fmt.Printf(" %v / %v allocated\n", len(n.Allocations), addrSpace)
		}
	}
}

//type IPv6Network struct{}

// Vlan represents vlan specification.
type Vlan struct {
	Model
	Description   string       `json:"description"`
	IPv4NetworkID uint         `json:"ipv4_network_id" gorm:"unique;not null"`
	IPv4Network   *IPv4Network `json:"ipv4_network,omitempty"`
	//IPv6NetworkID uint
	//IPv6Network IPv6Network
}

// IPv4Allocation represents allocated ip address specification.
type IPv4Allocation struct {
	Model
	Type        string `json:"type" gorm:"default:'generic'"` // generic or reserved
	Address     string `json:"address" gorm:"unique;not null" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//IPv4NetworkID uint         `json:"ipv4_network_id" gorm:"not null" validate:"required" sql:"type:integer REFERENCES ipv4_networks(id)"`
	IPv4NetworkID uint         `json:"ipv4_network_id" gorm:"not null" validate:"required"`
	IPv4Network   *IPv4Network `json:"ipv4_network,omitempty"`
}

// type IPv6Allocation struct {}

// GetNetwork returns net.IPNet instance address of IPv4Network.
func (n IPv4Network) GetNetwork() *net.IPNet {
	_, ipv4Net, _ := net.ParseCIDR(n.CIDR)
	return ipv4Net
}

// GetNetworkAddress returns network address string for the network.
func (n IPv4Network) GetNetworkAddress() string {
	ipNet := n.GetNetwork()
	return ipNet.IP.String()
}

// GetPrefixLength returns prefix size of the network.
func (n IPv4Network) GetPrefixLength() int {
	slice := strings.Split(n.CIDR, "/")
	len, _ := strconv.Atoi(slice[1])
	return len
}

// Contains return whether passed address is included the network or not.
func (n IPv4Network) Contains(addr string) bool {
	ip := net.ParseIP(addr)
	ipnet := n.GetNetwork()
	return ipnet.Contains(ip)
}

// GetAddressInteger returns Uint32 value represents its ip address.
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

// IPv4Networks is slice of IPv4Network
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

func (n IPv4Networks) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(n, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("ID	CIDR			Description\n")
		for _, network := range n {
			fmt.Printf("%2d	%-20s	%s\n", network.ID, network.CIDR, network.Description)
		}
	} else {
		if len(n[0].Subnetworks) > 0 {
			var writeNetworkTree func(networks IPv4Networks, depth int)
			writeNetworkTree = func(networks IPv4Networks, d int) {
				for _, n := range networks {
					fmt.Printf("%v%v\n", strings.Repeat("   ", d), n.CIDR)
					if len(n.Subnetworks) > 0 {
						writeNetworkTree(n.Subnetworks, d+1)
					}
				}
			}
			writeNetworkTree(n, 0)
		} else {
			fmt.Printf("%-20s	%s\n", "CIDR", "Description")
			for _, network := range n {
				fmt.Printf("%-20s	%s\n", network.CIDR, network.Description)
			}
		}
	}
}

// IPv4Allocations is slice of IPv4Allocation
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

func (a IPv4Allocations) Write(output string) {
	if output == "json" {
		jsonBytes, _ := json.MarshalIndent(a, "", "    ")
		fmt.Println(string(jsonBytes))
	} else if output == "wide" {
		fmt.Printf("%-5s %-15v  %-10v %-16v  %-v\n", "ID", "Address", "Type", "Name", "Description")
		for _, alloc := range a {
			fmt.Printf("%-5d %-15v  %-10v %-16v  %-v\n", alloc.ID, alloc.Address, alloc.Type, alloc.Name, alloc.Description)
		}
	} else {
		fmt.Printf("%-15v   %-15v   %-9v   %-16v   %-v\n", "Network", "Address", "Type", "Name", "Description")
		for _, alloc := range a {
			fmt.Printf("%-15v   %-15v   %-9v   %-16v   %-v\n", alloc.IPv4Network.CIDR, alloc.Address, alloc.Type, alloc.Name, alloc.Description)
		}
	}
}
