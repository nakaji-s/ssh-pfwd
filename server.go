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
	e.File("/", "ssh-pfwd/static/index.html")
	e.File("/vue.js", "ssh-pfwd/static/vue.js")
	e.File("/app.js", "ssh-pfwd/static/app.js")
	e.File("/axios.min.js", "ssh-pfwd/static/axios.min.js")
	e.File("/foundation.min.css", "ssh-pfwd/static/foundation.min.css")

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
