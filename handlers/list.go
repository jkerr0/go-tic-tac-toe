package handlers

import (
	"net/http"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/repository"
)

func Games(ctx Context, template string) error {
	games, err := repository.GetGames(ctx.Db)
	c := ctx.EchoCtx
	if err != nil {
		msg := "Could not retrieve games from database"
		c.Logger().Error(msg)
		return c.String(http.StatusInternalServerError, msg)
	}
	return c.Render(http.StatusOK, template, games)
}

func CreateGame(ctx Context) error {
	c := ctx.EchoCtx
	name := c.FormValue("name")
	err := repository.InsertGame(ctx.Db, name)
	if err != nil {
		msg := "could not create a game"
		c.Logger().Error(msg)
		return c.String(http.StatusInternalServerError, msg)
	}
	return Games(ctx, "games")
}

func DeleteGame(ctx Context) error {
	c := ctx.EchoCtx
	id, atoiErr := strconv.Atoi(c.Param("id"))
	if atoiErr != nil {
		msg := "Bad request id is not a number"
		c.Logger().Error(msg)
		return c.String(http.StatusBadRequest, msg)
	}
	return repository.DeleteGame(ctx.Db, id)
}
