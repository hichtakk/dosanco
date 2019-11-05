package main

import (
	"github.com/labstack/echo/v4"

	"github.com/hichikaw/dosanco/handler"
)

func setRoute(e *echo.Echo) {
	// routing network
	e.GET("/network", handler.GetAllNetwork)
	e.POST("/network", handler.CreateIPv4Network)
	e.GET("/network/:id", handler.GetIPv4Network)
	e.PUT("/network/:id", handler.UpdateIPv4Network)
	e.DELETE("/network/:id", handler.DeleteIPv4Network)

	e.GET("/ip/v4", handler.GetIPv4Allocations)
	e.POST("/ip/v4", handler.CreateIPv4Allocation)
	e.PUT("/ip/v4/:allocation_id", handler.UpdateIPv4Allocation)
	e.DELETE("/ip/v4/:allocation_id", handler.DeleteIPv4Allocation)

	e.GET("/vlan", handler.GetAllVlan)
	e.POST("/vlan", handler.CreateVlan)
	e.PUT("/vlan/:id", handler.UpdateVlan)
	e.DELETE("/vlan/:id", handler.DeleteVlan)

	// routing host
	e.GET("/host", handler.GetHosts)
	e.POST("/host", handler.CreateHost)
	e.GET("/host/:id", handler.GetHost)
	e.PUT("/host/:id", handler.UpdateHost)
	e.DELETE("/host/:id", handler.DeleteHost)
	//e.GET("/host/name/:hostname", handler.GetHostByName)

	// host groups
	e.GET("/host/group", handler.GetHostGroups)
	e.POST("/host/group", handler.CreateHostGroup)
	e.GET("/host/group/:id", handler.GetHostGroup)
	e.PUT("/host/group/:id", handler.UpdateHostGroup)
	e.DELETE("/host/group/:id", handler.DeleteHostGroup)

	// routing datacenter
	e.GET("/datacenter", handler.GetDataCenters)
	e.GET("/datacenter/:id", handler.GetDataCenter)
	e.POST("/datacenter", handler.CreateDataCenter)
	e.PUT("/datacenter/:id", handler.UpdateDataCenter)
	e.DELETE("/datacenter/:id", handler.DeleteDataCenter)
	//e.GET("/datacenter/name/:name", handler.GetDataCenterByName)

	// datacenter floor
	e.GET("/datacenter/floor", handler.GetAllDataCenterFloors)
	//e.GET("/datacenter/:id/floor", handler.GetDataCenterFloorsByDC)
	e.GET("/datacenter/floor/:id", handler.GetDataCenterFloor)
	//e.GET("/datacenter/floor/name/:name", handler.GetDataCenterFloorByName)
	e.POST("/datacenter/floor", handler.CreateDataCenterFloor)
	e.PUT("/datacenter/floor/:id", handler.UpdateDataCenterFloor)
	e.DELETE("/datacenter/floor/:id", handler.DeleteDataCenterFloor)

	// datacenter hall
	e.GET("/datacenter/hall", handler.GetDataCenterHalls)
	e.GET("/datacenter/hall/:id", handler.GetDataCenterHall)
	e.POST("/datacenter/hall", handler.CreateDataCenterHall)
	e.PUT("/datacenter/hall/:id", handler.UpdateDataCenterHall)
	e.DELETE("/datacenter/hall/:id", handler.DeleteDataCenterHall)

	// routing rack row
	e.GET("/datacenter/row", handler.GetRackRows)
	e.GET("/datacenter/row/:id", handler.GetRackRow)
	e.POST("/datacenter/row", handler.CreateRackRow)
	e.PUT("/datacenter/row/:id", handler.UpdateRackRow)
	e.DELETE("/datacenter/row/:id", handler.DeleteRackRow)

	// routing rack
	e.GET("/datacenter/rack", handler.GetRacks)
	e.GET("/datacenter/rack/:id", handler.GetRack)
	e.POST("/datacenter/rack", handler.CreateRack)
	e.PUT("/datacenter/rack/:id", handler.UpdateRack)
	e.DELETE("/datacenter/rack/:id", handler.DeleteRack)

	// datacenter power
	// routing UPS
	e.GET("/datacenter/ups", handler.GetUPSs)
	e.GET("/datacenter/ups/:id", handler.GetUPS)
	e.POST("/datacenter/ups", handler.CreateUPS)
	e.PUT("/datacenter/ups/:id", handler.UpdateUPS)
	e.DELETE("/datacenter/ups/:id", handler.DeleteUPS)

	// routing RowPDU
	e.GET("/datacenter/row-pdu", handler.GetRowPDUs)
	e.GET("/datacenter/row-pdu/:id", handler.GetRowPDU)
	e.POST("/datacenter/row-pdu", handler.CreateRowPDU)
	e.PUT("/datacenter/row-pdu/:id", handler.UpdateRowPDU)
	e.DELETE("/datacenter/row-pdu/:id", handler.DeleteRowPDU)

	// routing RackPDU
	e.GET("/datacenter/rack-pdu", handler.GetRackPDUs)
	e.GET("/datacenter/rack-pdu/:id", handler.GetRackPDU)
	e.POST("/datacenter/rack-pdu", handler.CreateRackPDU)
	e.PUT("/datacenter/rack-pdu/:id", handler.UpdateRackPDU)
	e.DELETE("/datacenter/rack-pdu/:id", handler.DeleteRackPDU)
}
