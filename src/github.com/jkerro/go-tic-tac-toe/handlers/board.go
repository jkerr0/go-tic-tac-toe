package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
)

func SelectSideAndGetBoard(ctx Context) error {
	c := ctx.EchoCtx
	side := c.Param("side")
	allowed := []string{"x", "o", "spectator"}
	sideCorrect := false
	for _, a := range allowed {
		if a == side {
			sideCorrect = true
		}
	}
	if !sideCorrect {
		return c.String(http.StatusBadRequest, "Invalid side")
	}
	db := ctx.Db
	session := GetSession(c)
	gameId, err := strconv.Atoi(session.Values["gameId"].(string))
	if err != nil {
		return c.String(http.StatusBadRequest, "Game id is required to be an integer")
	}

	if err = selectSide(db, session, gameId, side); err != nil {
		return c.String(http.StatusInternalServerError, "Cannot select side")
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

func selectSide(db *sqlx.DB, session *sessions.Session, gameId int, side string) error {
	var err error
	if side == "x" {
		err = repository.UpdateGameXUserId(db, gameId, session.Values["userId"].(int))
	} else if side == "o" {
		err = repository.UpdateGameOUserId(db, gameId, session.Values["userId"].(int))
	}
	if err != nil {
		return err
	}
	session.Values["side"] = side
	return nil
}
