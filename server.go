package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

type Server struct {
	Config Config
}

func (s Server) Start() {
	e := echo.New()
	e.File("/", "ssh-pfwd/static/index.html")
	e.File("/app.js", "ssh-pfwd/static/app.js")
	e.Static("/lib", "ssh-pfwd/static/lib")

	e.POST("/rule", func(c echo.Context) error {
		rule := new(Rule)
		if err := c.Bind(rule); err != nil {
			return err
		}
		rule.Id = uuid.NewV4().String()
		s.Config.AddRule(*rule)
		return c.JSON(http.StatusCreated, rule)
	})
	e.GET("/rules", func(c echo.Context) error {
		return c.JSON(http.StatusOK, s.Config.Rules)
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		if err := s.Config.DeleteRule(c.Param("id")); err != nil {
			return err
		}
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
	e.PUT("/rule/:id", func(c echo.Context) error {
		for i, _ := range s.Config.Rules {
			if s.Config.Rules[i].Id == c.Param("id") {
				if err := c.Bind(&s.Config.Rules[i]); err != nil {
					return err
				}
				return c.JSON(http.StatusOK, s.Config.Rules[i])
			}
		}
		return c.JSON(http.StatusNotFound, struct{}{})
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
