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
	gameId, err := strconv.Atoi(session.Values["gameId"].(string))
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
	SaveSession(session, c)
	return c.Render(http.StatusOK, "board", BoardData{b.Matrix(), gameId})
}
