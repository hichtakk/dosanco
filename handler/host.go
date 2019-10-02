package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	if result := db.Create(&host); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError(fmt.Sprintf("%v", result.Error)))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host created. ID: %d, Name: %s", host.ID, host.Name)))
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
	/*
		if result := db.Model(&h).Update("name", host.Name).Update("description", host.Description).Update("location", host.Location); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError("database error"))
		}
	*/
	if result := db.Model(&h).Update("name", host.Name).Update("description", host.Description); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("host updated. ID: %d, Name: %s, Description: %s", h.ID, h.Name, h.Description)))
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
