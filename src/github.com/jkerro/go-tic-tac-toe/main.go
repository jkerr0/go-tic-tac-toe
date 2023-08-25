package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/repository"
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
	db, err := sqlx.Connect("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	log.Println("Database connected")

	e := echo.New()
	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer
	e.GET("/", func(c echo.Context) error {
		games, err := repository.GetGames(db)
		if (err != nil) {
			panic(err)
		}
		return c.Render(http.StatusOK, "main", games)
	})

	e.POST("/games", func(c echo.Context) error {
		name := c.FormValue("name")
		repository.InsertGame(db, name)
		games, err := repository.GetGames(db)
		if (err != nil) {
			panic(err)
		}
		return c.Render(http.StatusOK, "games", games)
	})

	e.DELETE("/games/:id", func(c echo.Context) error {
		id, atoiErr := strconv.Atoi(c.Param("id"))
		if (atoiErr != nil) {
			return c.String(http.StatusBadRequest, "Bad request id is not a number")
		}
		return repository.DeleteGame(db, id)
	})


	e.Logger.Fatal(e.Start(":8080"))
}