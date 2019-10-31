package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

// GetDataCenters returns all datacenter information.
func GetDataCenters(c echo.Context) error {
	db := db.GetDB()
	dcs := []model.DataCenter{}
	name := c.QueryParam("name")
	if name != "" {
		db.Find(&dcs, "name=?", name)
		if len(dcs) == 0 {
			return c.JSON(http.StatusNotFound, returnError("dc not found"))
		}
	} else {
		db.Find(&dcs)
	}

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

// GetDataCenterFloorByName returns specified datacenter floor.
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
	existFloor := new(model.Floor)
	if result := db.Take(existFloor, "name=? AND data_center_id=?", floor.Name, floor.DataCenterID); result.RecordNotFound() != true {
		return c.JSON(http.StatusConflict, returnError(fmt.Sprintf("floor '%v' is already exist in %v", floor.Name, floor.DataCenterID)))
	}
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

// GetUPS returns datacenter floors.
func GetUPS(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse ups id error"))
	}
	ups := new(model.UPS)
	db := db.GetDB()
	if result := db.Take(&ups, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("ups not found"))
	}

	return c.JSON(http.StatusOK, ups)
}

// GetUPSs returns datacenter floors.
func GetUPSs(c echo.Context) error {
	db := db.GetDB()
	ups := model.UPSs{}
	dcName := c.QueryParam("dc")
	name := c.QueryParam("name")
	if dcName != "" {
		// get datacenter
		dc := model.DataCenter{}
		if result := db.Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// get ups
		if name != "" {
			db.Find(&ups, "data_center_id=? AND name=?", dc.ID, name)
		} else {
			db.Find(&ups, "data_center_id=?", dc.ID)
		}
	} else {
		if name != "" {
			db.Find(&ups, "name=?", name)
		} else {
			db.Find(&ups)
		}
	}

	return c.JSON(http.StatusOK, ups)
}

// GetRowPDUs returns datacenter floors.
func GetRowPDUs(c echo.Context) error {
	db := db.GetDB()
	pdu := model.RowPDUs{}
	dcName := c.QueryParam("dc")
	name := c.QueryParam("name")
	if dcName != "" {
		// get dc
		dc := model.DataCenter{}
		if result := db.Order("name").Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// get ups
		ups := model.UPSs{}
		if result := db.Order("name").Find(&ups, "data_center_id=?", dc.ID); result.Error != nil {
			return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("ups not found for dc '%v'", dcName)))
		}
		// get pdu
		for _, u := range ups {
			p := []model.RowPDU{}
			if name != "" {
				db.Order("name").Find(&p, "name=? AND primary_ups_id=?", name, u.ID)
			} else {
				db.Order("name").Find(&p, "primary_ups_id=?", u.ID)
			}
			pdu = append(pdu, p...)
		}
	} else {
		if name != "" {
			db.Order("name").Find(&pdu, "name=?", name)
		} else {
			db.Order("name").Find(&pdu)
		}
	}

	return c.JSON(http.StatusOK, pdu)
}

// GetRowPDU returns specified row-pdu information
func GetRowPDU(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse row-pdu id error"))
	}
	pdu := new(model.RowPDU)
	db := db.GetDB()
	if result := db.Take(&pdu, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("row-pdu not found"))
	}

	return c.JSON(http.StatusOK, pdu)
}

// GetRackPDU returns specified rack-pdu information
func GetRackPDU(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse rack-pdu id error"))
	}
	pdu := new(model.RackPDU)
	db := db.GetDB()
	if result := db.Take(&pdu, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("rack-pdu not found"))
	}

	return c.JSON(http.StatusOK, pdu)
}

// GetRackPDUs returns datacenter floors.
func GetRackPDUs(c echo.Context) error {
	db := db.GetDB()
	pdu := model.RackPDUs{}
	dcName := c.QueryParam("dc")
	upsName := c.QueryParam("ups")
	rowPduName := c.QueryParam("row-pdu")
	name := c.QueryParam("name")
	if dcName != "" {
		// get dc
		dc := model.DataCenter{}
		if result := db.Take(&dc, "name=?", dcName); result.Error != nil {
			return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("datacenter '%v' not found", dcName)))
		}
		// get ups
		ups := model.UPSs{}
		if upsName != "" {
			if result := db.Find(&ups, "name=? AND data_center_id=?", upsName, dc.ID); result.Error != nil {
				return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("ups '%v' not found for dc '%v'", upsName, dcName)))
			}
		} else {
			if result := db.Find(&ups, "data_center_id=?", dc.ID); result.Error != nil {
				return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("ups not found for dc '%v'", dcName)))
			}
		}
		// get dc pdu
		for _, u := range ups {
			dcPDUs := []model.RowPDU{}
			if rowPduName != "" {
				if result := db.Find(&dcPDUs, "name=? AND primary_ups_id=?", rowPduName, u.ID); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("row pdu '%v' not found", rowPduName)))
				}
			} else {
				if result := db.Find(&dcPDUs, "primary_ups_id=?", u.ID); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("dc pdu not found")))
				}
			}
			for _, dcPDU := range dcPDUs {
				p := []model.RackPDU{}
				if name != "" {
					db.Find(&p, "name=? AND primary_pdu_id=?", name, dcPDU.ID)
				} else {
					db.Find(&p, "primary_pdu_id=?", dcPDU.ID)
				}
				pdu = append(pdu, p...)
			}
		}
	} else {
		rowPDUs := []model.RowPDU{}
		if upsName != "" {
			ups := model.UPS{}
			if result := db.Take(&ups, "name=?", upsName); result.Error != nil {
				return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("ups '%v' not found", upsName)))
			}
			if rowPduName != "" {
				if result := db.Find(&rowPDUs, "name=? AND primary_ups_id=?", rowPduName, ups.ID); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("row-pdu '%v' not found under UPS '%v'", rowPduName, upsName)))
				}
			} else {
				if result := db.Find(&rowPDUs, "primary_ups_id=?", ups.ID); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError("row-pdu not found"))
				}
			}
			for _, rowPDU := range rowPDUs {
				p := model.RackPDUs{}
				if name != "" {
					db.Find(&p, "primary_pdu_id=? AND name=?", rowPDU.ID, name)
				} else {
					db.Find(&p, "primary_pdu_id=?", rowPDU.ID)
				}
				pdu = append(pdu, p...)
			}
		} else {
			if rowPduName != "" {
				if result := db.Find(&rowPDUs, "name=?", rowPduName); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError(fmt.Sprintf("row-pdu '%v' not found", rowPduName)))
				}
			} else {
				if result := db.Find(&rowPDUs); result.Error != nil {
					return c.JSON(http.StatusNotFound, returnError("row-pdu not found"))
				}
			}
			for _, rowPDU := range rowPDUs {
				p := model.RackPDUs{}
				if name != "" {
					db.Find(&p, "primary_pdu_id=? AND name=?", rowPDU.ID, name)
				} else {
					db.Find(&p, "primary_pdu_id=?", rowPDU.ID)
				}
				pdu = append(pdu, p...)
			}
		}
	}

	return c.JSON(http.StatusOK, pdu)
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
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("hall created. ID: %d, Name: %s", hall.ID, hall.Name)))
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

// GetRackRow returns specified datacenter rack row
func GetRackRow(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse row id error"))
	}
	row := new(model.RackRow)
	db := db.GetDB()
	if result := db.Take(&row, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("row not found"))
	}

	return c.JSON(http.StatusOK, row)
}

// GetRackRows returns datacenter rack rows.
func GetRackRows(c echo.Context) error {
	db := db.GetDB()
	rows := model.RackRows{}
	dcName := c.QueryParam("dc")
	floorName := c.QueryParam("floor")
	hallName := c.QueryParam("hall")
	name := c.QueryParam("name")

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
				if name != "" {
					db.Find(&rows, "name=? AND hall_id=?", name, hall.ID)
				} else {
					db.Find(&rows, "hall_id=?", hall.ID)
				}
			} else {
				halls := model.Halls{}
				db.Find(&halls, "floor_id=?", floor.ID)
				for _, hall := range halls {
					hallRows := model.RackRows{}
					if name != "" {
						db.Find(&hallRows, "name=? AND hall_id=?", name, hall.ID)
					} else {
						db.Find(&hallRows, "hall_id=?", hall.ID)
					}
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
					if name != "" {
						db.Find(&hallRows, "name=? AND hall_id=?", name, hall.ID)
					} else {
						db.Find(&hallRows, "hall_id=?", hall.ID)
					}
					rows = append(rows, hallRows...)
				}
			}
		}
	} else {
		if name != "" {
			db.Find(&rows, "name=?", name)
		} else {
			db.Find(&rows)
		}
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

// UpdateRackRow updates specified datacenter rack row information.
func UpdateRackRow(c echo.Context) error {
	row := new(model.RackRow)
	if err := c.Bind(row); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(rowID) != row.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched row ID between URI and request body."))
	}
	var r model.RackRow
	db := db.GetDB()
	if result := db.Take(&r, "id=?", rowID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("row not found on database."))
	}
	if result := db.Model(&r).Update("name", row.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack row updated. ID: %d, Name: %s", r.ID, row.Name)))
}

// DeleteRackRow deletes specified datacenter
func DeleteRackRow(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing row id error"))
	}
	db := db.GetDB()
	var row model.RackRow
	if result := db.Take(&row, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("row '%v' not found", id)))
	}
	if result := db.Delete(&row); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack row '%d' deleted", id)))
}

// GetRacks returns datacenter rack rows.
func GetRacks(c echo.Context) error {
	db := db.GetDB()
	racks := model.Racks{}
	dcName := c.QueryParam("dc")
	floorName := c.QueryParam("floor")
	hallName := c.QueryParam("hall")
	rowName := c.QueryParam("row")
	pduName := c.QueryParam("pdu")
	name := c.QueryParam("name")

	if dcName != "" {
		if pduName != "" {
			pdu := new(model.RowPDU)
			if result := db.Take(&pdu, "name=?", pduName); result.Error != nil {
				return c.JSON(http.StatusBadRequest, returnError("row-pdu not found"))
			}
			rackPDUs := []model.RackPDU{}
			db.Find(&rackPDUs, "primary_pdu_id=? OR secondary_pdu_id=?", pdu.ID, pdu.ID)
			if len(rackPDUs) == 0 {
				return c.JSON(http.StatusBadRequest, returnError("rack-pdu not found"))
			}
			hosts := []model.Host{}
			for _, rackPDU := range rackPDUs {
				host := new(model.Host)
				if result := db.Take(&host, "name=?", rackPDU.Name); result.Error == nil {
					hosts = append(hosts, *host)
				}
			}
			if len(hosts) == 0 {
				return c.JSON(http.StatusNotFound, returnError("rack-pdu host not found"))
			}
			rm := make(map[uint]struct{})
			for _, host := range hosts {
				rm[host.RackID] = struct{}{}
			}
			rs := []int{}
			for id := range rm {
				rs = append(rs, int(id))
			}
			sort.Sort(sort.IntSlice(rs))
			racks := []model.Rack{}
			for _, id := range rs {
				rack := model.Rack{}
				if result := db.Take(&rack, "id=?", id); result.Error == nil {
					racks = append(racks, rack)
				}
			}
			return c.JSON(http.StatusOK, racks)
		}
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
				if rowName != "" {
					row := model.RackRow{}
					if result := db.Take(&row, "name=? AND hall_id=?", rowName, hall.ID); result.Error != nil {
						return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("row '%v' not found", rowName)))
					}
					if name != "" {
						db.Find(&racks, "name=? AND row_id=?", name, row.ID)
					} else {
						db.Find(&racks, "row_id=?", row.ID)
					}
				} else {
					hallRows := model.RackRows{}
					db.Find(&hallRows, "hall_id=?", hall.ID)
					for _, row := range hallRows {
						rowRacks := model.Racks{}
						if name != "" {
							db.Find(&rowRacks, "name=? AND row_id=?", name, row.ID)
						} else {
							db.Find(&rowRacks, "row_id=?", row.ID)
						}
						racks = append(racks, rowRacks...)
					}
				}

			} else {
				halls := model.Halls{}
				db.Find(&halls, "floor_id=?", floor.ID)
				for _, hall := range halls {
					hallRows := model.RackRows{}
					db.Find(&hallRows, "hall_id=?", hall.ID)
					for _, row := range hallRows {
						rowRacks := model.Racks{}
						if name != "" {
							db.Find(&rowRacks, "name=? AND row_id=?", name, row.ID)
						} else {
							db.Find(&rowRacks, "row_id=?", row.ID)
						}
						racks = append(racks, rowRacks...)
					}
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
					for _, row := range hallRows {
						rowRacks := model.Racks{}
						if name != "" {
							db.Find(&rowRacks, "name=? AND row_id=?", name, row.ID)
						} else {
							db.Find(&rowRacks, "row_id=?", row.ID)
						}
						racks = append(racks, rowRacks...)
					}
				}
			}
		}
	} else {
		if name != "" {
			db.Find(&racks, "name=?", name)
		} else {
			db.Find(&racks)
		}
	}

	return c.JSON(http.StatusOK, racks)
}

// GetRack returns specified rack information.
func GetRack(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parse rack id error"))
	}
	rack := new(model.Rack)
	db := db.GetDB()
	if result := db.Take(&rack, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("rack not found"))
	}

	return c.JSON(http.StatusOK, rack)
}

// CreateRack creates a new rack row to specified data hall.
func CreateRack(c echo.Context) error {
	rack := new(model.Rack)
	if err := c.Bind(rack); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&rack); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack created. ID: %d, Name: %s", rack.ID, rack.Name)))
}

// UpdateRack updates specified datacenter rack row information.
func UpdateRack(c echo.Context) error {
	rack := new(model.Rack)
	if err := c.Bind(rack); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	rackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(rackID) != rack.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched row ID between URI and request body."))
	}
	var r model.Rack
	db := db.GetDB()
	if result := db.Take(&r, "id=?", rackID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("rack not found on database."))
	}
	if result := db.Model(&r).Update("name", rack.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack updated. ID: %d, Name: %s", r.ID, rack.Name)))
}

// DeleteRack deletes specified datacenter
func DeleteRack(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing row id error"))
	}
	db := db.GetDB()
	var rack model.Rack
	if result := db.Take(&rack, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("rack '%v' not found", id)))
	}
	if result := db.Delete(&rack); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack '%d' deleted", id)))
}

// CreateUPS creates a new rack row to specified data hall.
func CreateUPS(c echo.Context) error {
	ups := new(model.UPS)
	if err := c.Bind(ups); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&ups); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("ups created. ID: %d, Name: %s", ups.ID, ups.Name)))
}

// UpdateUPS updates specified UPS information.
func UpdateUPS(c echo.Context) error {
	ups := new(model.UPS)
	if err := c.Bind(ups); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	upsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(upsID) != ups.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched UPS ID between URI and request body."))
	}
	var u model.UPS
	db := db.GetDB()
	if result := db.Take(&u, "id=?", upsID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("ups not found on database."))
	}
	if result := db.Model(&u).Update("name", ups.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("ups updated. ID: %d, Name: %s", ups.ID, ups.Name)))
}

// UpdateRowPDU updates specified UPS information.
func UpdateRowPDU(c echo.Context) error {
	pdu := new(model.RowPDU)
	if err := c.Bind(pdu); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	pduID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(pduID) != pdu.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched PDU ID between URI and request body."))
	}
	var p model.RowPDU
	db := db.GetDB()
	if result := db.Take(&p, "id=?", pduID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("pdu not found on database."))
	}
	if result := db.Model(&p).Update("name", pdu.Name).Update("description", pdu.Description); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("pdu updated. ID: %d, Name: %s", pdu.ID, pdu.Name)))
}

// UpdateRackPDU updates specified UPS information.
func UpdateRackPDU(c echo.Context) error {
	pdu := new(model.RackPDU)
	if err := c.Bind(pdu); err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	pduID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError(err.Error()))
	}
	if uint(pduID) != pdu.ID {
		return c.JSON(http.StatusBadRequest, returnError("mismatched RACK PDU ID between URI and request body."))
	}
	var p model.RackPDU
	db := db.GetDB()
	if result := db.Take(&p, "id=?", pduID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("rack pdu not found on database."))
	}
	if result := db.Model(&p).Update("name", pdu.Name); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError("database write error."))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("pdu updated. ID: %d, Name: %s", pdu.ID, pdu.Name)))
}

// DeleteUPS deletes specified datacenter
func DeleteUPS(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing row id error"))
	}
	db := db.GetDB()
	var ups model.UPS
	if result := db.Take(&ups, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("ups '%v' not found", id)))
	}
	if result := db.Delete(&ups); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("ups '%d' deleted", id)))
}

// DeleteRowPDU deletes specified datacenter
func DeleteRowPDU(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing dc pdu id error"))
	}
	db := db.GetDB()
	var pdu model.RowPDU
	if result := db.Take(&pdu, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("pdu '%v' not found", id)))
	}
	if result := db.Delete(&pdu); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("pdu '%d' deleted", id)))
}

// DeleteRackPDU deletes specified datacenter
func DeleteRackPDU(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnError("parsing rack pdu id error"))
	}
	db := db.GetDB()
	var pdu model.RackPDU
	if result := db.Take(&pdu, "id=?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, returnError(fmt.Sprintf("rack pdu '%v' not found", id)))
	}
	if result := db.Delete(&pdu); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}

	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("pdu '%d' deleted", id)))
}

// CreateRowPDU creates a new rack row to specified data hall.
func CreateRowPDU(c echo.Context) error {
	pdu := new(model.RowPDU)
	if err := c.Bind(pdu); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()
	if result := db.Create(&pdu); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("pdu created. ID: %d, Name: %s", pdu.ID, pdu.Name)))
}

// CreateRackPDU creates a new rack row to specified data hall.
func CreateRackPDU(c echo.Context) error {
	pdu := new(model.RackPDU)
	if err := c.Bind(pdu); err != nil {
		return c.JSON(http.StatusBadRequest, returnError("received bad request: "+err.Error()))
	}
	db := db.GetDB()

	// validation
	existRackPDU := new(model.RackPDU)
	if result := db.Take(existRackPDU, "name=?", pdu.Name); result.RecordNotFound() != true {
		return c.JSON(http.StatusConflict, returnError(fmt.Sprintf("rack-pdu '%v' is already exist", pdu.Name)))
	}

	if result := db.Create(&pdu); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, returnError("database error"))
	}
	return c.JSON(http.StatusOK, returnMessage(fmt.Sprintf("rack pdu created. ID: %d, Name: %s", pdu.ID, pdu.Name)))
}
