package main

import (
	"net/http"
	"os"

	"time"

	"io"

	"github.com/GeertJohan/go.rice"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/satori/go.uuid"
)

type Server struct {
	Config Config
}

func (s Server) Start() {
	e := echo.New()

	// log rotate settings
	logf, err := rotatelogs.New(
		"./access_log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("./access_log"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
		return
	}
	// logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// log to stdout and file
		Output: io.MultiWriter(os.Stdout, logf),
	}))
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

	// Set ResponseHeader
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set no-cache
			c.Response().Header().Set("Cache-Control", "no-cache")
			return next(c)
		}
	})

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
		err := s.Config.AddRule(*rule)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusCreated, rule)
	})
	e.GET("/rules", func(c echo.Context) error {
		rules, err := s.Config.GetRules()
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusOK, rules)
	})
	e.DELETE("/rule/:id", func(c echo.Context) error {
		if err := s.Config.DeleteRule(c.Param("id")); err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusOK, struct{}{})
	})
	e.GET("/rule/:id", func(c echo.Context) error {
		rule, err := s.Config.GetRule(c.Param("id"))
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusOK, rule)
	})
	e.PUT("/rule/:id", func(c echo.Context) error {
		updatedRule, err := s.Config.UpdateRule(c.Param("id"), c)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusOK, updatedRule)
	})

	// TODO CRUD for key

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}

func handleError(c echo.Context, err error) error {
	switch err {
	case gorm.ErrRecordNotFound, ErrIdNotFound:
		return c.JSON(http.StatusNotFound, struct{}{})
	default:
		return err
	}
}
