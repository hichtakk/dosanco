package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/hichikaw/dosanco/config"
	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/handler"
	"github.com/hichikaw/dosanco/model"
)

type Validator struct {
	validator *validator.Validate
}

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
	e.Validator = &Validator{validator: validator.New()}
	
	// initialize logger middleware

	// route requests
	e.GET("/network", handler.GetNetwork)
	e.POST("/network", func(c echo.Context) error {
		network := new(model.IPv4Network)
		if err := c.Bind(network); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "received bad request. " + err.Error()})
		}
		if err := c.Validate(network); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "request validation failed. " + err.Error()})
		}
		if err := handler.CreateNetwork(network); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		} else {
			return c.JSON(http.StatusOK, network)
		}
	})
	e.GET("/network/:id", handler.GetNetwork)
	e.PUT("/network/:id", handler.UpdateNetwork)
	e.DELETE("/network/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := handler.DeleteNetwork(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		} else {
			return c.JSON(http.StatusOK, map[string]string{"message": "network deleted"})
		}
	})

	// Start dosanco server
	e.Logger.Fatal(e.Start(":8080"))
}