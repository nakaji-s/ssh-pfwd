package main

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/satori/go.uuid"
)

type Server struct {
	Config Config
}

func (s Server) Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if len(reqBody) > 0 {
			e.Logger.Debug(string(reqBody))
		}
	}))

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin" {
			return true, nil
		}
		return false, nil
	}))

	e.File("/", "ssh-pfwd/static/index.html")
	e.File("/app.js", "ssh-pfwd/static/app.js")
	assetHandler := http.FileServer(rice.MustFindBox("static/lib").HTTPBox())
	e.GET("/lib/*", echo.WrapHandler(http.StripPrefix("/lib/", assetHandler)))

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
		updatedRule, err := s.Config.UpdateRule(c.Param("id"), c)
		if err != nil {
			return c.JSON(http.StatusNotFound, struct{}{})
		}
		return c.JSON(http.StatusOK, updatedRule)
	})

	// TODO CRUD for key

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
