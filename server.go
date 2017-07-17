package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type Server struct {
	Config Config
}

func (s Server) Start() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.PUT("/rule", func(c echo.Context) error {
		return c.String(http.StatusOK, "Add rule!")
	})
	e.GET("/rule", func(c echo.Context) error {
		return c.String(http.StatusOK, "GET rules!")
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "Delete rule! :"+c.Param("id"))
	})
	e.GET("/rule/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "GET rule! :"+c.Param("id"))
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
