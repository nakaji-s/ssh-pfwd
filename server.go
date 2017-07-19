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
		return c.JSON(http.StatusOK, s.Config.Rules)
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		s.Config.DeleteRule(c.Param("id"))
		return c.JSON(http.StatusOK, struct{}{})
	})
	e.GET("/rule/:id", func(c echo.Context) error {
		for _, rule := range s.Config.Rules {
			if rule.Id == c.Param("id") {
				return c.JSON(http.StatusOK, rule)
			}
		}
		return c.JSON(http.StatusNotFound, struct{}{})
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
