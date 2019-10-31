package handler

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

// GetAllNetwork returns all network information.
func GetAllNetwork(c echo.Context) error {
	db := db.GetDB()
	networks := model.IPv4Networks{}
	tree, _ := strconv.ParseBool(c.QueryParam("tree"))
	rfc, _ := strconv.ParseBool(c.QueryParam("show-rfc-reserved"))
	cidr := c.QueryParam("cidr")
	if tree == true {
		// 0: all, other: specified nubmer of depth
		depth, _ := strconv.Atoi(c.QueryParam("depth"))
		if depth <= 0 {
			depth = -1
		}
		root := model.IPv4Network{}
		db.Take(&root, "id=1")
		subnets := getSubnetworks(root.ID, uint(depth), uint(0), rfc)
		root.Subnetworks = *subnets
		networks = append(networks, root)
	} else {
		// return flat network list
		if cidr != "" {
			db.Find(&networks, "c_id_r=?", cidr)
		} else {
			db.Find(&networks)
		}
		sort.Sort(networks)
	}

	return c.JSON(http.StatusOK, networks)
}

// GetIPv4Network returns a specified network information.
func GetIPv4Network(c echo.Context) error {
	var network model.IPv4Network
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse network id error"))
	}
	db := db.GetDB()
	if result := db.Take(&network, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("network not found"))
	}
	subnets := getSubnetworks(network.ID, 1, uint(0), false)
	if len(*subnets) > 0 {
		network.Subnetworks = *subnets
	}
	allocs := getIPAllocations(network.ID)
	if len(*allocs) > 0 {
		network.Allocations = *allocs
	}
	return c.JSON(http.StatusOK, network)
}

// GetIPv4NetworkByCIDR returns a specified CIDR network information.
func GetIPv4NetworkByCIDR(c echo.Context) error {
	var network model.IPv4Network
	cidr := strings.Replace(c.Param("cidr"), "-", "/", 1)
	db := db.GetDB()
	if result := db.Take(&network, "c_id_r=?", cidr); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("network not found"))
	}
	subnets := getSubnetworks(network.ID, 1, uint(0), false)
	if len(*subnets) > 0 {
		network.Subnetworks = *subnets
	}
	allocs := getIPAllocations(network.ID)
	if len(*allocs) > 0 {
		network.Allocations = *allocs
	}
	return c.JSON(http.StatusOK, network)
}

// CreateIPv4Network creates a new network with given json data.
func CreateIPv4Network(c echo.Context) error {
	network := new(model.IPv4Network)
	if err := c.Bind(network); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request. "+err.Error()))
	}
	if err := c.Validate(network); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("request validation failed. "+err.Error()))
	}
	ipv4Addr, _, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		return err
	}
	if ipv4Addr.String() != network.GetNetworkAddress() {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("'%v' is not network address", network.CIDR)))
	}

	db := db.GetDB()
	var root model.IPv4Network
	db.Take(&root, "c_id_r=?", "0.0.0.0/0")
	supernet := GetSupernetwork(&root, network)

	// ensure the network is not overwrapped other subnets.
	subnets := []model.IPv4Network{}
	db.Where(&model.IPv4Network{SupernetworkID: supernet.ID}).Find(&subnets)
	for _, s := range subnets {
		if s.GetPrefixLength() > network.GetPrefixLength() {
			if network.Contains(s.GetNetworkAddress()) {
				return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("requested network '%v' is overwrapping with network '%v'", network.CIDR, s.CIDR)))
			}
		}
		n := s.GetNetwork()
		if n.Contains(ipv4Addr) {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("requested network '%v' is overwrapping with network '%v'", network.CIDR, n)))
		}
	}

	network.SupernetworkID = supernet.ID
	if result := db.Create(&network); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError(result.Error.Error()))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("network created. ID: %d,  CIDR: %s,  Description: %s", network.ID, network.CIDR, network.Description)))
}

// GetSupernetwork returns supernetwork for network passed to argument.
func GetSupernetwork(seed *model.IPv4Network, nw *model.IPv4Network) *model.IPv4Network {
	subnets := getSubnetworks(seed.ID, uint(1), uint(0), true)
	for _, subnet := range *subnets {
		if subnet.GetNetworkAddress() == nw.GetNetworkAddress() {
			if subnet.GetPrefixLength() > nw.GetPrefixLength() {
				return seed
			}
		}
		if subnet.Contains(nw.GetNetworkAddress()) {
			return GetSupernetwork(&subnet, nw)
		}
	}

	return seed
}

// UpdateIPv4Network updates only description for specified network.
func UpdateIPv4Network(c echo.Context) error {
	network := new(model.IPv4Network)
	if err := c.Bind(network); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	network.ID = uint(id)
	db := db.GetDB()
	var net model.IPv4Network
	if result := db.Take(&net, "id=?", network.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("network not found"))
	}
	if result := db.Model(&net).Update("description", network.Description); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("network updated. ID: %d,  CIDR: %s,  Description: %s", net.ID, net.CIDR, network.Description)))
}

// DeleteIPv4Network deletes specified network
func DeleteIPv4Network(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if id == 1 {
		return c.JSON(http.StatusBadRequest, returnError("cannot delete root network"))
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
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("network has subnets [%v]", strings.Join(subnetList, ", "))))
	}
	// ensure the network does not have ip allocations
	allocations := []model.IPv4Allocation{}
	db.Where(&model.IPv4Allocation{IPv4NetworkID: network.ID}).Find(&allocations)
	if len(allocations) > 0 {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("network has ip allocations")))
	}
	// ensure the network does not have vlan
	vlan := model.Vlan{}
	if db.Where(&model.Vlan{IPv4NetworkID: network.ID}).Take(&vlan).RecordNotFound() == false {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("network has vlan %v", vlan.ID)))
	}

	db.Unscoped().Delete(&network)

	return c.JSON(http.StatusOK, returnMessage("network deleted"))
}

// GetIPv4Allocations returns ip allocation data for specified network.
func GetIPv4Allocations(c echo.Context) error {
	nid, err := strconv.Atoi(c.Param("network_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	addr := getIPAllocations(uint(nid))

	return c.JSON(http.StatusOK, addr)
}

// CreateIPv4Allocation creates a new ipv4 allocation to specified network.
func CreateIPv4Allocation(c echo.Context) error {
	addr := new(model.IPv4Allocation)
	if err := c.Bind(addr); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request. "+err.Error()))
	}
	if err := c.Validate(addr); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("request validation failed. "+err.Error()))
	}
	if (addr.Type == "generic") && (addr.Name == "") {
		return c.JSON(http.StatusBadRequest, returnError("hostname is required for type: generic"))
	}

	db := db.GetDB()
	network := new(model.IPv4Network)
	if result := db.Take(network, "id=?", addr.IPv4NetworkID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("network not found"))
	}
	if network.Contains(addr.Address) != true {
		return c.JSON(http.StatusBadRequest, returnError("requested address is not an address of specified network"))
	}
	subnets := getSubnetworks(network.ID, uint(1), uint(0), true)
	if len(*subnets) != 0 {
		return c.JSON(http.StatusBadRequest, returnError("requested network is subnetted"))
	}

	if result := db.Create(addr); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database request failed. "+result.Error.Error()))
	}

	return c.JSON(http.StatusOK, returnMessage("ip allocation created"))
}

// DeleteIPv4Allocation deletes specified ipv4 allocation.
func DeleteIPv4Allocation(c echo.Context) error {
	allocID, err := strconv.Atoi(c.Param("allocation_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	addr := new(model.IPv4Allocation)
	db := db.GetDB()
	if result := db.Delete(addr, "id=?", allocID); result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, returnMessage("ip allocation deleted"))
}

// UpdateIPv4Allocation updates only description of specified ipv4 allocation.
func UpdateIPv4Allocation(c echo.Context) error {
	allocID, err := strconv.Atoi(c.Param("allocation_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	addr := new(model.IPv4Allocation)
	reqAddr := new(model.IPv4Allocation)
	if err := c.Bind(reqAddr); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request. "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Take(addr, "id=?", allocID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("target not found"))
	}
	if addr.Type == "generic" && reqAddr.Name == "" {
		return c.JSON(http.StatusBadRequest, returnError("hostname is required for generic ip allocation"))
	}
	if addr.Type == "reserved" && reqAddr.Name != "" {
		return c.JSON(http.StatusBadRequest, returnError("cannot be set hostname to reserved ip allocation"))
	}
	if result := db.Model(addr).Update("description", reqAddr.Description).Update("name", reqAddr.Name); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage("allocation update"))
}

// GetHostIPv4Allocations returns ipv4 allocations associated with specified hostname.
func GetHostIPv4Allocations(c echo.Context) error {
	hostname := c.Param("hostname")
	addr := []model.IPv4Allocation{}
	d := db.GetDB()
	d.Find(&addr, "name=?", hostname)

	return c.JSON(http.StatusOK, addr)
}

// GetIPv4AllocationByAddress returns ipv4 allocation associated with specified address.
func GetIPv4AllocationByAddress(c echo.Context) error {
	ipv4 := c.Param("address")
	addr := model.IPv4Allocation{}
	d := db.GetDB()
	d.Find(&addr, "address=?", ipv4)

	return c.JSON(http.StatusOK, addr)
}

func getSubnetworks(id uint, depth uint, step uint, rfc bool) *model.IPv4Networks {
	var subnetworks model.IPv4Networks
	subnetworkList := []model.IPv4Network{}
	db := db.GetDB()
	if rfc == true {
		db.Where("supernetwork_id=?", id).Find(&subnetworkList)
	} else {
		db.Where("supernetwork_id=? and reserved=?", id, false).Find(&subnetworkList)
	}
	if (len(subnetworkList) > 0) && (step < depth) {
		for _, sn := range subnetworkList {
			gsn := getSubnetworks(sn.ID, depth, step+1, rfc)
			sn.Subnetworks = append(sn.Subnetworks, *gsn...)
			subnetworks = append(subnetworks, sn)
		}
	}
	sort.Sort(subnetworks)

	return &subnetworks
}

func getIPAllocations(id uint) *model.IPv4Allocations {
	var allocs model.IPv4Allocations
	db := db.GetDB()
	db.Where(&model.IPv4Allocation{IPv4NetworkID: id}).Find(&allocs)
	sort.Sort(allocs)
	return &allocs
}

func getIPv4AllocationsByHostname(name string) *model.IPv4Allocations {
	var allocs model.IPv4Allocations
	db := db.GetDB()
	db.Where(&model.IPv4Allocation{Name: name}).Find(&allocs)
	sort.Sort(allocs)
	return &allocs
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
		return c.JSON(http.StatusBadRequest, returnError("received bad request. "+err.Error()))
	}
	if vlan.ID > 4094 {
		return c.JSON(http.StatusBadRequest, returnError("request validation failed. Vlan ID should be the range of 1 - 4094."))
	}
	db := db.GetDB()
	var network model.IPv4Network
	if result := db.Take(&network, "id=?", vlan.IPv4NetworkID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("network not found"))
	}
	if result := db.Create(&vlan); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError(fmt.Sprintf("%v", result.Error)))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("vlan created. ID: %d, Description: %s", vlan.ID, vlan.Description)))
}

// UpdateVlan updates only description for specified vlan
func UpdateVlan(c echo.Context) error {
	vlan := new(model.Vlan)
	if err := c.Bind(vlan); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	vlan.ID = uint(id)
	db := db.GetDB()
	var v model.Vlan
	if result := db.Take(&v, "id=?", vlan.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("vlan not found"))
	}
	if result := db.Model(&v).Update("description", vlan.Description); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("vlan updated. ID: %d, Description: %s", v.ID, vlan.Description)))
}

// DeleteVlan deletes specified vlan
func DeleteVlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	db := db.GetDB()
	var vlan model.Vlan
	if result := db.Take(&vlan, "id=?", id); result.Error != nil {
		return fmt.Errorf("vlan '%v' not found", id)
	}
	db.Unscoped().Delete(&vlan)

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("vlan %d deleted", id)))
}
