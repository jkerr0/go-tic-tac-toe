package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jkerro/go-tic-tac-toe/handlers"
	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
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
	repository.CreateDatabase(db)

	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer
	context := func(c echo.Context) handlers.Context {
		return handlers.Context{EchoCtx: c, Db: db}
	}

	e.GET("/", func(c echo.Context) error {
		return handlers.Games(context(c), "main")
	})

	e.GET("/games", func(c echo.Context) error {
		return handlers.Games(context(c), "games")
	})

	e.POST("/games", func(c echo.Context) error {
		return handlers.CreateGame(context(c))
	})

	e.DELETE("/games/:id", func(c echo.Context) error {
		return handlers.DeleteGame(context(c))
	})

	e.GET("/board/:gameId", func(c echo.Context) error {
		gameId, err := strconv.Atoi(c.Param("gameId"))
		if err != nil {
			c.String(http.StatusBadRequest, "Game id is required to be an integer")
		}
		moves, err := repository.GetMoves(db, gameId)
		if err != nil {
			c.String(http.StatusInternalServerError, "Database error")
		}
		b := logic.GetBoard(moves)
		type BoardData struct {
			Board  [][]logic.BoardElement
			GameId int
		}
		return c.Render(http.StatusOK, "board", BoardData{b.Matrix(), gameId})
	})

	channels := handlers.NewChannelPool()

	e.GET("/ws/:gameId", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		gameId, err := strconv.Atoi(c.Param("gameId"))
		if err != nil {
			c.String(http.StatusBadRequest, "Game id is required to be an integer")
		}

		channel := channels.Join(ws, gameId)

		defer func() {
			channels.Leave(ws, gameId)
			ws.Close()
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			} else {
				handlers.HandleMoveMessage(context(c), msg, gameId, channel)
			}
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":8080"))
}
