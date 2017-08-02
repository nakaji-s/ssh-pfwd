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
		return c.JSON(http.StatusOK, s.Config.GetRules())
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		if err := s.Config.DeleteRule(c.Param("id")); err != nil {
			return c.JSON(http.StatusNotFound, struct{}{})
		}
		return c.JSON(http.StatusOK, struct{}{})
	})
	e.GET("/rule/:id", func(c echo.Context) error {
		rule, err := s.Config.GetRule(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusNotFound, struct{}{})
		}
		return c.JSON(http.StatusOK, rule)
	})
	e.PUT("/rule/:id", func(c echo.Context) error {
		rule := new(Rule)
		if err := c.Bind(rule); err != nil {
			return err
		}
		updatedRule, err := s.Config.UpdateRule(c.Param("id"), *rule)
		if err != nil {
			return c.JSON(http.StatusNotFound, struct{}{})
		}
		return c.JSON(http.StatusOK, updatedRule)
	})

	// TODO CRUD for key

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
