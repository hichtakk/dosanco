package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	flag "github.com/spf13/pflag"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/hichikaw/dosanco/config"
	"github.com/hichikaw/dosanco/db"
	"github.com/hichikaw/dosanco/model"
)

const version = "v0.1.0"

var revision = ""

// Validator echo middleware
type Validator struct {
	validator *validator.Validate
}

// Validate function
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func showVersion() {
	fmt.Printf("dosanco apiserver \U0001f434, version: %s, revision: %s\n", version, revision)
}

func main() {
	// parse flags
	var configfile string
	var showversion bool
	flag.StringVarP(&configfile, "config", "c", "/etc/dosanco/config.toml", "configuration file path")
	flag.BoolVarP(&showversion, "version", "v", false, "show dosanco apiserver version")
	port := flag.UintP("port", "p", 15187, "dosanco-apiserver listening port")
	flag.Parse()
	if showversion == true {
		showVersion()
		os.Exit(0)
	}

	// read configuration
	conf, err := config.NewConfig(configfile)
	if err != nil {
		panic(err.Error())
	}

	// initialize database
	db.Init(conf)

	// initialize echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = false
	e.Validator = &Validator{validator: validator.New()}

	// initialize logger middleware
	e.Use(middleware.Logger())
	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.Version{Version: version, Revision: revision})
	})
	setRoute(e)

	// Start dosanco server
	listenPort := strconv.Itoa(int(*port))
	e.Logger.Fatal(e.Start(":" + listenPort))
}
