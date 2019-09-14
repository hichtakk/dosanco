package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/hichikaw/dosanco/config"
	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/handler"
)

// Validator echo middleware
type Validator struct {
	validator *validator.Validate
}

// Validate function
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func main() {
	// read configuration
	conf, err := config.NewConfig("./config.toml")
	if err != nil {
		panic(err.Error())
	}

	// initialize database
	db.Init(conf)

	// initialize echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &Validator{validator: validator.New()}

	// initialize logger middleware
	e.Use(middleware.Logger())

	/*
		Routing Definition
	*/
	// routing network
	e.GET("/network", handler.GetAllNetwork)
	e.POST("/network", handler.CreateIPv4Network)
	e.GET("/network/:id", handler.GetIPv4Network)
	e.PUT("/network/:id", handler.UpdateIPv4Network)
	e.DELETE("/network/:id", handler.DeleteIPv4Network)
	e.GET("/network/cidr/:cidr", handler.GetIPv4NetworkByCIDR)

	e.POST("/ip/v4", handler.CreateIPv4Allocation)
	e.PUT("/ip/v4/:allocation_id", handler.UpdateIPv4Allocation)
	e.DELETE("/ip/v4/:allocation_id", handler.DeleteIPv4Allocation)
	e.GET("/ip/v4/network/:network_id", handler.GetIPv4Allocations)
	e.GET("/ip/v4/host/:hostname", handler.GetHostIPv4Allocations)
	e.GET("/ip/v4/addr/:address", handler.GetIPv4AllocationByAddress)

	e.GET("/vlan", handler.GetAllVlan)
	e.POST("/vlan", handler.CreateVlan)
	e.PUT("/vlan/:id", handler.UpdateVlan)
	e.DELETE("/vlan/:id", handler.DeleteVlan)

	// routing host
	e.POST("/host", handler.CreateHost)
	e.GET("/host/:id", handler.GetHost)
	e.PUT("/host/:id", handler.UpdateHost)
	e.DELETE("/host/:id", handler.DeleteHost)
	e.GET("/host/name/:hostname", handler.GetHostByName)

	// routing datacenter
	e.GET("/datacenter", handler.GetAllDataCenters)
	e.GET("/datacenter/:id", handler.GetDataCenter)
	e.POST("/datacenter", handler.CreateDataCenter)
	e.PUT("/datacenter/:id", handler.UpdateDataCenter)
	e.DELETE("/datacenter/:id", handler.DeleteDataCenter)
	e.GET("/datacenter/name/:name", handler.GetDataCenterByName)

	// datacenter floor
	e.GET("/datacenter/floor", handler.GetAllDataCenterFloors)
	e.GET("/datacenter/:id/floor", handler.GetDataCenterFloorsByDC)
	e.GET("/datacenter/floor/:id", handler.GetDataCenterFloor)
	e.GET("/datacenter/floor/name/:name", handler.GetDataCenterFloorByName)
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
	e.POST("/datacenter/row", handler.CreateRackRow)

	// Start dosanco server
	e.Logger.Fatal(e.Start(":8080"))
}
