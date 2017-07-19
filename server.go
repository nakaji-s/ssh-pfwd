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
		rule := new(Rule)
		if err := c.Bind(rule); err != nil {
			return err
		}
		s.Config.AddRule(*rule)
		return c.JSON(http.StatusOK, rule)
	})
	e.GET("/rules", func(c echo.Context) error {
		//return c.String(http.StatusOK, "GET rules!")
		return c.JSON(http.StatusOK, s.Config.Rules)
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "Delete rule! :"+c.Param("id"))
	})
	e.GET("/rule/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "GET rule! :"+c.Param("id"))
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
