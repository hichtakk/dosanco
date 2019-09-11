package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

// GetAllDataCenters returns all datacenter information.
func GetAllDataCenters(c echo.Context) error {
	db := db.GetDB()
	dcs := []model.DataCenter{}
	db.Find(&dcs)

	return c.JSON(http.StatusOK, dcs)
}

// GetDataCenter returns specified datacenter information.
func GetDataCenter(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse dc id error"))
	}
	dc := new(model.DataCenter)
	db := db.GetDB()
	if result := db.Take(&dc, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("dc not found"))
	}

	flr := new([]model.Floor)
	db.Find(&flr, "data_center_id=?", dc.ID)
	dc.Floors = *flr

	return c.JSON(http.StatusOK, dc)
}

// GetDataCenterByName returns specified host information.
func GetDataCenterByName(c echo.Context) error {
	dc := new(model.DataCenter)
	db := db.GetDB()
	if result := db.Take(&dc, "name=?", c.Param("name")); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("dc not found"))
	}

	return c.JSON(http.StatusOK, dc)
}

// CreateDataCenter creates a new data center.
func CreateDataCenter(c echo.Context) error {
	dc := new(model.DataCenter)
	if err := c.Bind(dc); err != nil {
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()}, "    ")
	}
	db := db.GetDB()
	if result := db.Create(&dc); result.Error != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("%v", result.Error)}, "    ")
	}
	return c.JSONPretty(http.StatusOK, map[string]string{"message": fmt.Sprintf("data center created. ID: %d, Address: %s", dc.ID, dc.Address)}, "    ")
}

// UpdateDataCenter updates address for specified datacenter
func UpdateDataCenter(c echo.Context) error {
	dc := new(model.DataCenter)
	if err := c.Bind(dc); err != nil {
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "    ")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "    ")
	}
	dc.ID = uint(id)
	db := db.GetDB()
	var d model.DataCenter
	if result := db.Take(&d, "id=?", dc.ID); result.Error != nil {
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "datacenter not found"}, "    ")
	}
	if result := db.Model(&d).Update("address", dc.Address); result.Error != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "database error"}, "    ")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"message": fmt.Sprintf("datacenter updated. ID: %d, Address: %s", d.ID, dc.Address)}, "    ")
}

// DeleteDataCenter deletes specified datacenter
func DeleteDataCenter(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "    ")
	}
	db := db.GetDB()
	var dc model.DataCenter
	if result := db.Take(&dc, "id=?", id); result.Error != nil {
		return fmt.Errorf("datacenter '%v' not found", id)
	}
	db.Delete(&dc)

	return c.JSONPretty(http.StatusOK, map[string]string{"message": fmt.Sprintf("datacenter %d deleted", id)}, "    ")
}

// GetAllDataCenterFloors returns all of datacenter floors.
func GetAllDataCenterFloors(c echo.Context) error {
	db := db.GetDB()
	flrs := []model.Floor{}
	db.Find(&flrs)

	return c.JSON(http.StatusOK, flrs)
}

// GetDataCenterFloorsByDC returns datacenter floors of specified datacenter.
func GetDataCenterFloorsByDC(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse dc id error"))
	}
	dc := new(model.DataCenter)
	db := db.GetDB()
	if result := db.Take(&dc, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("dc not found"))
	}

	flrs := new(model.Floors)
	if result := db.Find(&flrs, "data_center_id=?", dc.ID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database error"))
	}

	return c.JSON(http.StatusOK, flrs)
}

// GetDataCenterFloor returns specified datacenter floor.
func GetDataCenterFloor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse floor id error"))
	}
	flr := new(model.Floor)
	db := db.GetDB()
	if result := db.Take(&flr, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("floor not found"))
	}

	return c.JSON(http.StatusOK, flr)
}
