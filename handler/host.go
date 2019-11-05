package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

// GetHost returns specified host information.
func GetHost(c echo.Context) error {
	host := new(model.Host)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse host id error"))
	}
	db := db.GetDB()
	if result := db.Take(&host, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("host not found"))
	}
	allocs := getIPv4AllocationsByHostname(host.Name)
	if len(*allocs) > 0 {
		host.IPv4Allocations = *allocs
	}

	return c.JSON(http.StatusOK, host)
}

// GetHosts returns specified host information.
func GetHosts(c echo.Context) error {
	db := db.GetDB()
	groupName := c.QueryParam("group")
	location := c.QueryParam("location")

	hosts := new(model.Hosts)
	if groupName == "" && location == "" {
		return c.JSON(http.StatusBadRequest, returnError("query 'group' or 'location' is required"))
	}
	group := new(model.HostGroup)
	if groupName != "" {
		if result := db.Take(group, "name=?", groupName); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("group '%v' not found", groupName)))
		}
	}
	rack := new(model.Rack)
	if location != "" {
		location, err := url.QueryUnescape(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, returnError(err.Error()))
		}
		locSlice := strings.Split(location, "/")
		if len(locSlice) != 5 {
			return c.JSON(http.StatusBadRequest, returnError("invalid location format. use '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'"))
		}
		dcName := locSlice[0]
		floorName := locSlice[1]
		hallName := locSlice[2]
		rowName := locSlice[3]
		rackName := locSlice[4]
		// datacenter
		dc := new(model.DataCenter)
		if result := db.Take(dc, "name=?", dcName); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// floor
		floor := new(model.Floor)
		if result := db.Take(floor, "name=? AND data_center_id=?", floorName, dc.ID); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("floor '%v' not found", floorName)))
		}
		// hall
		hall := new(model.Hall)
		if result := db.Take(hall, "name=? AND floor_id=?", hallName, floor.ID); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("hall '%v' not found", hallName)))
		}
		// row
		row := new(model.RackRow)
		if result := db.Take(row, "name=? AND hall_id=?", rowName, hall.ID); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("row '%v' not found", rowName)))
		}
		// rack
		if result := db.Take(rack, "name=? AND row_id=?", rackName, row.ID); result.RecordNotFound() == true {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("rack '%v' not found", rackName)))
		}
	}
	if group.ID != 0 && rack.ID != 0 {
		db.Find(hosts, "group_id=? AND rack_id=?", group.ID, rack.ID)
	} else if group.ID != 0 && rack.ID == 0 {
		db.Find(hosts, "group_id=?", group.ID)
	} else if group.ID == 0 && rack.ID != 0 {
		db.Find(hosts, "rack_id=?", rack.ID)
	}

	return c.JSON(http.StatusOK, hosts)
}

// GetHostGroup returns specified host information.
func GetHostGroup(c echo.Context) error {
	group := new(model.HostGroup)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse host id error"))
	}
	db := db.GetDB()
	if result := db.Take(&group, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("host group not found"))
	}

	return c.JSON(http.StatusOK, group)
}

// GetHostGroups returns host group information.
func GetHostGroups(c echo.Context) error {
	db := db.GetDB()
	groups := model.HostGroups{}
	name := c.QueryParam("name")
	if name != "" {
		db.Find(&groups, "name=?", name)
	} else {
		db.Find(&groups)
	}

	return c.JSON(http.StatusOK, groups)
}

// GetHostByName returns specified host information.
func GetHostByName(c echo.Context) error {
	host := new(model.Host)
	db := db.GetDB()
	if result := db.Take(&host, "name=?", c.Param("hostname")); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("host not found"))
	}
	allocs := getIPv4AllocationsByHostname(host.Name)
	if len(*allocs) > 0 {
		host.IPv4Allocations = *allocs
	}

	return c.JSON(http.StatusOK, host)
}

// CreateHost creates a new host.
func CreateHost(c echo.Context) error {
	host := new(model.Host)
	if err := c.Bind(host); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request. "+err.Error()))
	}
	db := db.GetDB()
	h := new(model.Host)
	if db.Take(&h, "name=?", host.Name).RecordNotFound() != true {
		return c.JSON(http.StatusConflict, returnError(fmt.Sprintf("host '%v' is already exist", host.Name)))
	}
	if result := db.Create(&host); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError(fmt.Sprintf("%v", result.Error)))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host created. ID: %d, Name: %s", host.ID, host.Name)))
}

// CreateHostGroup creates a new host group.
func CreateHostGroup(c echo.Context) error {
	group := new(model.HostGroup)
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	exist := new(model.HostGroup)
	if result := db.Take(&exist, "name=?", group.Name); result.Error == nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("host group '%v' is already exist", group.Name)))
	}
	if result := db.Create(&group); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host group created. ID: %d, Name: %s", group.ID, group.Name)))
}

// UpdateHost updates information of specified host.
func UpdateHost(c echo.Context) error {
	host := new(model.Host)
	if err := c.Bind(host); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	host.ID = uint(id)
	db := db.GetDB()
	var h model.Host
	if result := db.Take(&h, "id=?", host.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("host not found"))
	}
	if result := db.Model(&h).Update("name", host.Name).Update("description", host.Description); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database error"))
	}
	if result := db.Model(&h).Update("group_id", host.GroupID).Update("rack_id", host.RackID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host updated. ID: %d, Name: %s, Description: %s", h.ID, h.Name, h.Description)))
}

// UpdateHostGroup updates information of specified host group.
func UpdateHostGroup(c echo.Context) error {
	group := new(model.HostGroup)
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	group.ID = uint(id)
	db := db.GetDB()
	var g model.HostGroup
	if notFound := db.Take(&g, "name=? AND id!=?", group.Name, group.ID).RecordNotFound(); notFound == false {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("host group '%v' is already exists", group.Name)))
	}
	if result := db.Take(&g, "id=?", group.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("host group not found"))
	}
	if result := db.Model(&g).Update("name", group.Name).Update("description", group.Description); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host group updated. ID: %d, Name: %s, Description: %s", g.ID, g.Name, g.Description)))
}

// DeleteHost deletes specified host.
func DeleteHost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	db := db.GetDB()
	var host model.Host
	if result := db.Take(&host, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("host %d not found", id)))
	}

	// ensure the host does not have allocated ip
	var alloc []model.IPv4Allocation
	db.Find(&alloc, "name=?", host.Name)
	if len(alloc) > 0 {
		return c.JSON(http.StatusBadRequest, returnError("host has ip allocations"))
	}

	if result := db.Delete(&host, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host %d is deleted", host.ID)))
}

// DeleteHostGroup deletes specified host.
func DeleteHostGroup(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	db := db.GetDB()
	var group model.HostGroup
	if result := db.Take(&group, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("host group %d not found", id)))
	}

	// ensure the host group does not have host
	var hosts []model.Host
	db.Find(&hosts, "group_id=?", group.ID)
	if len(hosts) > 0 {
		return c.JSON(http.StatusBadRequest, returnError("group has hosts"))
	}

	if result := db.Delete(&group, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host group %d is deleted", group.ID)))
}
