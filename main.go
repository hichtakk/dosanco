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
	db.Init(conf.DB)

	// initialize echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &Validator{validator: validator.New()}

	// initialize logger middleware
	e.Use(middleware.Logger())

	// route requests
	e.GET("/network", handler.GetAllNetwork)
	e.POST("/network", handler.CreateIPv4Network)
	e.GET("/network/:id", handler.GetIPv4Network)
	e.PUT("/network/:id", handler.UpdateIPv4Network)
	e.DELETE("/network/:id", handler.DeleteIPv4Network)

	e.POST("/ipam", handler.CreateIPv4Allocation)
	e.PUT("/ipam/:allocation_id", handler.UpdateIPv4Allocation)
	e.DELETE("/ipam/:allocation_id", handler.DeleteIPv4Allocation)
	e.GET("/ipam/network/:network_id", handler.GetIPv4Allocations)
	e.GET("/ipam/host/:hostname", handler.GetHostIPv4Allocations)

	e.GET("/vlan", handler.GetAllVlan)
	e.POST("/vlan", handler.CreateVlan)
	e.PUT("/vlan/:id", handler.UpdateVlan)
	e.DELETE("/vlan/:id", handler.DeleteVlan)

	e.GET("/datacenter", handler.GetAllDataCenters)
	e.POST("/datacenter", handler.CreateDataCenter)
	e.PUT("/datacenter/:id", handler.UpdateDataCenter)
	e.DELETE("/datacenter/:id", handler.DeleteDataCenter)

	// Start dosanco server
	e.Logger.Fatal(e.Start(":8080"))
}
