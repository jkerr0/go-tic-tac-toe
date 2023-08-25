package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

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
	repository.InsertGame(db, "newgame")
	games, err := repository.GetGames(db)
	for _, game := range games {
		log.Println(game.Name)
	}

	e := echo.New()
	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main", games)
	})


	e.Logger.Fatal(e.Start(":8080"))
}