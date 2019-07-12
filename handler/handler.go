package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo"
	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

func GetNetwork(c echo.Context) error {
	db := db.GetDB()
	networks := []model.IPv4Network{}
	db.Find(&networks)

	return c.JSON(http.StatusOK, networks)
}

func CreateNetwork(network *model.IPv4Network) error {
	db := db.GetDB()

	var supernet model.IPv4Network
	db.First(&supernet, "id=?", network.SupernetworkID)
	fmt.Printf("supernet: %v\n", supernet)
	supernetCIDR := supernet.GetNetwork()

	ipv4Addr, _, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("request: %v\n", network)
	if !supernetCIDR.Contains(ipv4Addr) {
		return fmt.Errorf("%v is out of %v\n", ipv4Addr, supernet.CIDR)
	}
	if network.GetPrefixLength() <= supernet.GetPrefixLength() {
		return fmt.Errorf("network '%v' is larger than supernetwork '%v'", network.CIDR, supernet.CIDR)
	}

	// ensure the network is not overwrapped other subnets.
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: network.SupernetworkID}).Find(&subnets)
	for _, s := range subnets {
		fmt.Printf("subnet: %v\n", s)
		n := s.GetNetwork()
		if n.Contains(ipv4Addr) {
			return fmt.Errorf("requested network '%v' is overwrapping with network '%v'", network.CIDR, n)
		}
	}

	db.Create(&network)

	return nil
}

func UpdateNetwork(c echo.Context) error {
	return nil
}

func DeleteNetwork(id int) error {
	db := db.GetDB()
	var network model.IPv4Network
	db.First(&network, "id=?", id)
	if network.ID != 0 {
		fmt.Println(network)
	} else {
		return fmt.Errorf("network not found")
	}
	// ensure the network does not have subnetworks
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: network.ID}).Find(&subnets)
	fmt.Println(subnets)
	if len(subnets) > 0 {
		return fmt.Errorf("network has subnets %v", subnets)
	}
	db.Delete(&network)

	return nil
}