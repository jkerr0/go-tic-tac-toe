package handlers

import (
	"net/http"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Games(db *sqlx.DB, template string, c echo.Context) error {
	games, err := repository.GetGames(db)
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, template, games)
}

func CreateGame(db *sqlx.DB, c echo.Context) error {
	name := c.FormValue("name")
	repository.InsertGame(db, name)
	games, err := repository.GetGames(db)
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, "games", games)
}

func DeleteGame(db *sqlx.DB, c echo.Context) error {
	id, atoiErr := strconv.Atoi(c.Param("id"))
	if atoiErr != nil {
		return c.String(http.StatusBadRequest, "Bad request id is not a number")
	}
	return repository.DeleteGame(db, id)
}
