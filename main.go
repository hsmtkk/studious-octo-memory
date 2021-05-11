package main

import (
	"io"
	"log"
	"text/template"

	"github.com/hsmtkk/studious-octo-memory/handle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t

	handler, err := handle.New("data.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/", handler.Index)

	e.Logger.Fatal(e.Start(":8000"))
}
