package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

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

	return c.JSONPretty(http.StatusOK, networks, "    ")
}

// GetNetwork returns a specified network
func GetIPv4Network(c echo.Context) error {
	var network model.IPv4Network
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ErrorResponse{Error: Error{Message: err.Error()}}, "    ")
	}
	db := db.GetDB()
	if result := db.Take(&network, "id=?", id); result.Error != nil {
		return c.JSONPretty(http.StatusBadRequest, returnBusinessError("network not found"), "    ")
	}
	return c.JSONPretty(http.StatusOK, network, "    ")
}

// CreateNetwork creates a new network with given json data
func CreateIPv4Network(c echo.Context) error {
	network := new(model.IPv4Network)
	if err := c.Bind(network); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()})
	}
	if err := c.Validate(network); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "request validation failed. " + err.Error()})
	}
	db := db.GetDB()
	var supernet model.IPv4Network
	db.Take(&supernet, "id=?", network.SupernetworkID)
	supernetCIDR := supernet.GetNetwork()
	ipv4Addr, _, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		return err
	}
	if !supernetCIDR.Contains(ipv4Addr) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("%v is out of %v", ipv4Addr, supernet.CIDR)})
	}
	if network.GetPrefixLength() <= supernet.GetPrefixLength() {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("network '%v' is larger than supernetwork '%v'", network.CIDR, supernet.CIDR)})
	}
	// ensure the network is not overwrapped other subnets.
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: network.SupernetworkID}).Find(&subnets)
	for _, s := range subnets {
		n := s.GetNetwork()
		if n.Contains(ipv4Addr) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("requested network '%v' is overwrapping with network '%v'", network.CIDR, n)})
		}
	}
	if result := db.Create(&network); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("%v", result.Error)})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("network created. ID: %d,  CIDR: %s,  Description: %s", network.ID, network.CIDR, network.Description)})
}

// UpdateNetwork updates only description for specified network
func UpdateIPv4Network(c echo.Context) error {
	network := new(model.IPv4Network)
	if err := c.Bind(network); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	network.ID = uint(id)
	db := db.GetDB()
	var net model.IPv4Network
	if result := db.Take(&net, "id=?", network.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "network not found"})
	}
	if result := db.Model(&net).Update("description", network.Description); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "database error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("network updated. ID: %d,  CIDR: %s,  Description: %s", net.ID, net.CIDR, network.Description)})
}

// DeleteNetwork deletes specified network
func DeleteIPv4Network(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if id == 1 {
		return fmt.Errorf("cannot delete root network")
	}
	db := db.GetDB()
	var network model.IPv4Network
	if result := db.Take(&network, "id=?", id); result.Error != nil {
		return fmt.Errorf("network '%v' not found", id)
	}
	// ensure the network does not have subnetworks
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: network.ID}).Find(&subnets)
	if len(subnets) > 0 {
		var subnetList []string
		for _, s := range subnets {
			subnetList = append(subnetList, s.CIDR)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("network has subnets [%v]", strings.Join(subnetList, ", "))})
	}
	db.Delete(&network)

	return c.JSON(http.StatusOK, map[string]string{"message": "network deleted"})
}

func GetIPv4Allocations(c echo.Context) error {
	nid, err := strconv.Atoi(c.Param("network_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	addr := []model.IPv4Allocation{}
	d := db.GetDB()
	d.Find(&addr, "ipv4_network_id=?", nid)

	return c.JSON(http.StatusOK, addr)
}

func CreateIPv4Allocation(c echo.Context) error {
	addr := new(model.IPv4Allocation)
	if err := c.Bind(addr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()})
	}
	if err := c.Validate(addr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "request validation failed. " + err.Error()})
	}
	db := db.GetDB()
	if result := db.Create(addr); result.Error != nil {
		//return result.Error
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "database request failed. " + result.Error.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ip allocation created"})
}

func DeleteIPv4Allocation(c echo.Context) error {
	allocId, err := strconv.Atoi(c.Param("allocation_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	addr := new(model.IPv4Allocation)
	db := db.GetDB()
	if result := db.Delete(addr, "id=?", allocId); result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ip allocation deleted"})
}

func UpdateIPv4Allocation(c echo.Context) error {
	allocId, err := strconv.Atoi(c.Param("allocation_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	addr := new(model.IPv4Allocation)
	reqAddr := new(model.IPv4Allocation)
	if err := c.Bind(reqAddr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()})
	}

	db := db.GetDB()
	if result := db.Take(addr, "id=?", allocId); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "target not found"})
	}

	if result := db.Model(addr).Update("description", reqAddr.Description); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "database error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "allocation update"})
}

func GetHostIPv4Allocations(c echo.Context) error {
	hostname := c.Param("hostname")
	addr := []model.IPv4Allocation{}
	d := db.GetDB()
	d.Find(&addr, "name=?", hostname)

	return c.JSON(http.StatusOK, addr)
}

func getSubnetworks(id uint, depth uint, step uint) *[]model.IPv4Network {
	subnetworks := []model.IPv4Network{}
	subnetworkList := []model.IPv4Network{}
	db := db.GetDB()
	db.Where(&model.IPv4Network{SupernetworkID: id}).Find(&subnetworkList)

	if (len(subnetworkList) > 0) && (step < depth) {
		for _, sn := range subnetworkList {
			gsn := getSubnetworks(sn.ID, depth, step+1)
			sn.Subnetworks = append(sn.Subnetworks, *gsn...)
			subnetworks = append(subnetworks, sn)
		}
	}

	return &subnetworks
}

// GetAllVlan returns all vlans
func GetAllVlan(c echo.Context) error {
	db := db.GetDB()
	vlans := []model.Vlan{}
	db.Find(&vlans)

	return c.JSONPretty(http.StatusOK, vlans, "    ")
}

// CreateVlan creates a new vlan with given json data
func CreateVlan(c echo.Context) error {
	vlan := new(model.Vlan)
	if err := c.Bind(vlan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()})
	}
	if err := c.Validate(vlan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "request validation failed. " + err.Error()})
	}
	if vlan.ID > 4094 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "request validation failed. Vlan ID should be the range of 1 - 4094."})
	}
	db := db.GetDB()
	if result := db.Create(&vlan); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("%v", result.Error)})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("vlan created. ID: %d, Description: %s", vlan.ID, vlan.Description)})
}

// UpdateVlan updates only description for specified vlan
func UpdateVlan(c echo.Context) error {
	vlan := new(model.Vlan)
	if err := c.Bind(vlan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	vlan.ID = uint(id)
	db := db.GetDB()
	var v model.Vlan
	if result := db.Take(&v, "id=?", vlan.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "vlan not found"})
	}
	if result := db.Model(&v).Update("description", vlan.Description); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "database error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("vlan updated. ID: %d, Description: %s", v.ID, vlan.Description)})
}

// DeleteVlan deletes specified vlan
func DeleteVlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	db := db.GetDB()
	var vlan model.Vlan
	if result := db.Take(&vlan, "id=?", id); result.Error != nil {
		return fmt.Errorf("vlan '%v' not found", id)
	}
	db.Delete(&vlan)

	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("vlan %d deleted\n", id)})
}
