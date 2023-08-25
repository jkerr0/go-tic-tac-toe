package main

import (
	"log"
	"net/http"

	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "modernc.org/sqlite"
)
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
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	e.Logger.Fatal(e.Start(":8080"))
}