package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

// GetAllNetwork returns all networks
func GetAllNetwork(c echo.Context) error {
	db := db.GetDB()
	networks := []model.IPv4Network{}

	tree, _ := strconv.ParseBool(c.QueryParam("tree"))
	if tree == true {
		// 0: all, other: specified nubmer of depth
		depth, _ := strconv.Atoi(c.QueryParam("depth"))
		if depth <= 0 {
			depth = -1
		}
		root := model.IPv4Network{}
		db.Take(&root, "id=1")
		subnets := getSubnetworks(root.ID, uint(depth), uint(0))
		root.Subnetworks = *subnets
		networks = append(networks, root)
	} else {
		// return flat network list
		db.Find(&networks)
	}

	return c.JSON(http.StatusOK, networks)
}

// GetNetwork returns a specified network
func GetNetwork(id int, network *model.IPv4Network) error {
	db := db.GetDB()
	if result := db.First(network, "id=?", id); result.Error != nil {
		return fmt.Errorf("network '%v' not found", id)
	}
	fmt.Println(network)

	return nil
}

// CreateNetwork creates a new network with given json data
func CreateNetwork(network *model.IPv4Network) error {
	db := db.GetDB()

	var supernet model.IPv4Network
	db.First(&supernet, "id=?", network.SupernetworkID)
	supernetCIDR := supernet.GetNetwork()

	ipv4Addr, _, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		return err
	}
	if !supernetCIDR.Contains(ipv4Addr) {
		return fmt.Errorf("%v is out of %v", ipv4Addr, supernet.CIDR)
	}
	if network.GetPrefixLength() <= supernet.GetPrefixLength() {
		return fmt.Errorf("network '%v' is larger than supernetwork '%v'", network.CIDR, supernet.CIDR)
	}

	// ensure the network is not overwrapped other subnets.
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: network.SupernetworkID}).Find(&subnets)
	for _, s := range subnets {
		n := s.GetNetwork()
		if n.Contains(ipv4Addr) {
			return fmt.Errorf("requested network '%v' is overwrapping with network '%v'", network.CIDR, n)
		}
	}
	if result := db.Create(&network); result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateNetwork updates only description for specified network
func UpdateNetwork(network *model.IPv4Network) (*model.IPv4Network, error) {
	db := db.GetDB()

	var net model.IPv4Network
	if result := db.Take(&net, "id=?", network.ID); result.Error != nil {
		return &model.IPv4Network{}, fmt.Errorf("network '%v' not found", network.ID)
	}
	if result := db.Model(&net).Update("description", network.Description); result.Error != nil {
		return &model.IPv4Network{}, result.Error
	}

	return &net, nil
}

// DeleteNetwork deletes specified network
func DeleteNetwork(id int) error {
	if id == 1 {
		return fmt.Errorf("cannot delete root network")
	}
	db := db.GetDB()
	var network model.IPv4Network
	if result := db.First(&network, "id=?", id); result.Error != nil {
		return fmt.Errorf("network '%v' not found", id)
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

func getSubnetworks(id uint, depth uint, step uint) *[]model.IPv4Network {
	subnetworks := []model.IPv4Network{}
	subnetworkList := []model.IPv4Network{}
	db := db.GetDB()
	db.Where(&model.IPv4Network{SupernetworkID: id}).Find(&subnetworkList)

	if (len(subnetworkList) > 0) && (step < depth) {
		//fmt.Printf("find more than one subnet for %v\n", id)
		for _, sn := range subnetworkList {
			gsn := getSubnetworks(sn.ID, depth, step+1)
			sn.Subnetworks = append(sn.Subnetworks, *gsn...)
			subnetworks = append(subnetworks, sn)
		}
	}

	return &subnetworks
}
