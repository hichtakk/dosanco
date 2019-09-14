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

	halls := new([]model.Hall)
	if db.Find(&halls, "floor_id=?", flr.ID).RecordNotFound() == false {
		flr.Halls = *halls
	}

	return c.JSON(http.StatusOK, flr)
}

// GetDataCenterFloor returns specified datacenter floor.
func GetDataCenterFloorByName(c echo.Context) error {
	name, _ := strconv.Atoi(c.Param("name"))
	flr := new(model.Floor)
	db := db.GetDB()
	if result := db.Take(&flr, "name=?", name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("floor not found"))
	}

	return c.JSON(http.StatusOK, flr)
}

// CreateDataCenterFloor creates a new floor to specified datacenter.
func CreateDataCenterFloor(c echo.Context) error {
	floor := new(model.Floor)
	if err := c.Bind(floor); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&floor); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("floor created. ID: %d, Address: %s", floor.ID, floor.Name)))
}

// UpdateDataCenterFloor updates specified datacenter floor information.
func UpdateDataCenterFloor(c echo.Context) error {
	floor := new(model.Floor)
	if err := c.Bind(floor); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	floorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(floorID) != floor.ID {
		return c.JSON(http.StatusBadRequest, returnError("floor ID specified in URI and request body are mismatched."))
	}
	var flr model.Floor
	db := db.GetDB()
	if result := db.Take(&flr, "id=?", floorID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("floor not found on database."))
	}
	if result := db.Model(&flr).Update("name", floor.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("datacenter floor updated. ID: %d, Name: %s", flr.ID, floor.Name)))
}

// DeleteDataCenterFloor deletes specified datacenter
func DeleteDataCenterFloor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing floor id error"))
	}
	db := db.GetDB()
	var floor model.Floor
	if result := db.Take(&floor, "id=?", id); result.Error != nil {
		return fmt.Errorf("floor '%v' not found", id)
	}
	db.Delete(&floor)

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("datacenter floor %d deleted", id)))
}

// GetAllDataCenterHalls returns all of datacenter floors.
func GetAllDataCenterHalls(c echo.Context) error {
	db := db.GetDB()
	halls := model.Halls{}
	db.Find(&halls)

	return c.JSON(http.StatusOK, halls)
}

// GetDataCenterHalls returns datacenter floors.
func GetDataCenterHalls(c echo.Context) error {
	db := db.GetDB()
	halls := model.Halls{}
	dcName := c.QueryParam("dc")
	floorName := c.QueryParam("floor")
	name := c.QueryParam("name")
	if (dcName == "") && (floorName != "") {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("'dc' parameter is required when 'floor' parameter is used.")))
	} else if (dcName != "") && (floorName == "") {
		// get datacenter
		dc := model.DataCenter{}
		if result := db.Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// get floor
		floors := []model.Floor{}
		db.Find(&floors, "data_center_id=?", dc.ID)
		for _, floor := range floors {
			// get hall for each floor
			h := model.Halls{}
			if name != "" {
				db.Find(&h, "floor_id=? AND name=?", floor.ID, name)
				halls = append(halls, h...)
			} else {
				db.Find(&h, "floor_id=?", floor.ID)
				halls = append(halls, h...)
			}
		}
	} else if (dcName != "") && (floorName != "") {
		// get datacenter
		dc := model.DataCenter{}
		if result := db.Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// get floor
		floor := model.Floor{}
		if result := db.Take(&floor, "name=? AND data_center_id=?", floorName, dc.ID); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("floor '%v' not found", floorName)))
		}
		if name != "" {
			db.Find(&halls, "floor_id=? AND name=?", floor.ID, name)
		} else {
			db.Find(&halls, "floor_id=?", floor.ID)
		}
	} else {
		if name != "" {
			db.Find(&halls, "name=?", name)
		} else {
			db.Find(&halls)
		}
	}

	return c.JSON(http.StatusOK, halls)
}

// GetDataCenterHall returns all of datacenter floors.
func GetDataCenterHall(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing hall id error"))
	}
	db := db.GetDB()
	hall := model.Hall{}
	if result := db.Find(&hall, "id=?", id); result.Error != nil {
		return fmt.Errorf("hall '%v' not found", id)
	}

	return c.JSON(http.StatusOK, hall)
}

// CreateDataCenterHall creates a new floor to specified datacenter.
func CreateDataCenterHall(c echo.Context) error {
	hall := new(model.Hall)
	if err := c.Bind(hall); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&hall); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("hall created. ID: %d, Name: %s, Type: %s", hall.ID, hall.Name, hall.Type)))
}

// UpdateDataCenterHall updates specified datacenter hall information.
func UpdateDataCenterHall(c echo.Context) error {
	hall := new(model.Hall)
	if err := c.Bind(hall); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	hallID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(hallID) != hall.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched hall ID between URI and request body."))
	}
	var h model.Hall
	db := db.GetDB()
	if result := db.Take(&h, "id=?", hallID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("hall not found on database."))
	}
	if result := db.Model(&h).Update("name", hall.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("datacenter hall updated. ID: %d, Name: %s", h.ID, hall.Name)))
}

// DeleteDataCenterHall deletes specified datacenter
func DeleteDataCenterHall(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing hall id error"))
	}
	db := db.GetDB()
	var hall model.Hall
	if result := db.Take(&hall, "id=?", id); result.Error != nil {
		return fmt.Errorf("hall '%v' not found", id)
	}
	if result := db.Delete(&hall); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("datacenter hall %d deleted", id)))
}

// GetRackRows returns datacenter rack rows.
func GetRackRows(c echo.Context) error {
	db := db.GetDB()
	rows := model.RackRows{}
	dcName := c.QueryParam("dc")
	floorName := c.QueryParam("floor")
	hallName := c.QueryParam("hall")

	if dcName != "" {
		dc := model.DataCenter{}
		if result := db.Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		if floorName != "" {
			floor := model.Floor{}
			if result := db.Take(&floor, "name=? AND data_center_id=?", floorName, dc.ID); result.Error != nil {
				return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("floor '%v' not found", floorName)))
			}
			if hallName != "" {
				hall := model.Hall{}
				if result := db.Take(&hall, "name=? AND floor_id=?", hallName, floor.ID); result.Error != nil {
					return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("hall '%v' not found", hallName)))
				}
				db.Find(&rows, "hall_id=?", hall.ID)
			} else {
				halls := model.Halls{}
				db.Find(&halls, "floor_id=?", floor.ID)
				for _, hall := range halls {
					hallRows := model.RackRows{}
					db.Find(&hallRows, "hall_id=?", hall.ID)
					rows = append(rows, hallRows...)
				}
			}
		} else {
			floors := model.Floors{}
			db.Find(&floors, "data_center_id=?", dc.ID)
			for _, floor := range floors {
				halls := model.Halls{}
				db.Find(&halls, "floor_id=?", floor.ID)
				for _, hall := range halls {
					hallRows := model.RackRows{}
					db.Find(&hallRows, "hall_id=?", hall.ID)
					rows = append(rows, hallRows...)
				}
			}
		}
	} else {
		db.Find(&rows)
	}

	return c.JSON(http.StatusOK, rows)
}

// CreateRackRow creates a new rack row to specified data hall.
func CreateRackRow(c echo.Context) error {
	row := new(model.RackRow)
	if err := c.Bind(row); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&row); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack row created. ID: %d, Name: %s", row.ID, row.Name)))
}
