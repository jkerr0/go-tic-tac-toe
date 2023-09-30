package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo-contrib/session"

	"github.com/antonlindstrom/pgstore"

	"github.com/gorilla/websocket"
	"github.com/jkerro/go-tic-tac-toe/handlers"
	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	if err != nil {
		e.Logger.Fatal("Could not load env file")
	}
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbPort := os.Getenv("POSTGRES_PORT")
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPassword, dbPort, dbName)
	e.Logger.Debug(connectionString)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		e.Logger.Fatal("Could not connect to the database. ", err)
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
				c.Logger().Error("Cannot create user", err)
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot create user"))
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
			return c.String(http.StatusBadRequest, "Game id is required to be an integer")
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
	sessionKey, defined := os.LookupEnv("SESSION_KEY")
	if !defined {
		e.Logger.Fatal("session key not defined")
	}
	sessionKeyByte := []byte(sessionKey)

	e.Static("/style", "/public/style")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "/",
		Browse: false,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	store, err := pgstore.NewPGStore(connectionString, sessionKeyByte)
	if err != nil {
		e.Logger.Fatal("Could not create postgres session store", err)
	}
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store:   store,
		Skipper: middleware.DefaultSkipper,
	}))
	e.Logger.Fatal(e.Start(":8080"))
}
