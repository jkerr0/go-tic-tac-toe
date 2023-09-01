package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jkerro/go-tic-tac-toe/handlers"
	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var (
	upgrader = websocket.Upgrader{}
)

func main() {
	e := echo.New()
	db, err := sqlx.Connect("sqlite", "test.db")
	if err != nil {
		log.Println("Could not connect to the database")
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

	channel := handlers.OpenChannel()

	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		channel.Join(ws)

		defer ws.Close()
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			} else {
				channel.Broadcast(fmt.Sprintf("<div id=\"receive\">%s</div>", string(msg)))
			}
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":8080"))
}
