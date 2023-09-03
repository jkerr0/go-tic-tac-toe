package handlers

import (
	"net/http"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
)

func GetBoard(ctx Context) error {
	c := ctx.EchoCtx
	db := ctx.Db
	session := GetSession(c)
	c.Logger().Info("joined with side ", session.Values["side"])
	session.Values["side"] = string(logic.X)
	SaveSession(session, c)
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
}
