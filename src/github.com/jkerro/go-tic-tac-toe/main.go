package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/jkerro/go-tic-tac-toe/handlers"
	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "modernc.org/sqlite"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	db, err := sqlx.Connect("sqlite", "test.db")
	if err != nil {
		e.Logger.Info("Could not connect to the database", err)
	}

	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer
	e.GET("/", func(c echo.Context) error {
		return handlers.Games(db, "main", c)
	})

	e.GET("/games", func(c echo.Context) error {
		return handlers.Games(db, "games", c)
	})

	e.POST("/games", func(c echo.Context) error {
		return handlers.CreateGame(db, c)
	})

	e.DELETE("/games/:id", func(c echo.Context) error {
		return handlers.DeleteGame(db, c)
	})

	e.GET("/board", func(c echo.Context) error {
		b := logic.NewBoard()
		return c.Render(http.StatusOK, "board", b.Matrix())
	})

	e.Logger.Fatal(e.Start(":8080"))
}
