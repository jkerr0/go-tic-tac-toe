package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/gommon/log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/gorilla/websocket"
	"github.com/jkerro/go-tic-tac-toe/handlers"
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
		e.Logger.Fatal("Could not connect to the database")
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
		s := handlers.GetSession(c)
		if s.Values["userId"] == nil {
			userId, err := repository.CreateUser(db)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Cannot create user")
			}
			s.Values["userId"] = userId
		}
		handlers.SaveSession(s, c)
		return handlers.Games(context(c), "main")
	})

	e.GET("/games", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.Games(context(c), "games")
		})
	})

	e.GET("/games-list-elements", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.Games(context(c), "games-list-elements")
		})
	})

	e.POST("/games", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.CreateGame(context(c))
		})
	})

	e.POST("/select-side/:side", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.SelectSideAndGetBoard(context(c))
		})
	})

	e.DELETE("/games/:id", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.DeleteGame(context(c))
		})
	})

	e.GET("/board/:gameId", func(c echo.Context) error {
		return handlers.HtmxOnly(c, func(c echo.Context) error {
			return handlers.CheckSideAndGetBoard(context(c))
		})
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

		side := handlers.GetSession(c).Values[fmt.Sprintf("side-%d", gameId)].(string)
		channel := channels.Join(ws, gameId)

		defer func() {
			channels.Leave(ws, gameId)
			ws.Close()
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			handlers.HandleMoveMessage(context(c), msg, gameId, channel, side)
		}

	})
	sessionKey := "secret"
	// sessionKey, defined := os.LookupEnv("SESSION_KEY")
	// if !defined {
	// 	e.Logger.Fatal("session key not defined")
	// }
	sessionKeyByte := []byte(sessionKey)

	e.Static("/style", "/public/style")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "/",
		Browse: false,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store:   sessions.NewFilesystemStore("/session"),
		Skipper: middleware.DefaultSkipper,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore(sessionKeyByte)))
	e.Logger.SetLevel(log.INFO)
	e.Logger.Fatal(e.Start(":8080"))
}
