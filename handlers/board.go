package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
)

type SideSelectorData struct {
	XSelected       bool
	OSelected       bool
	AlreadySelected bool
}

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
	gameId := session.Values["gameId"].(int)
	game, err := repository.GetGame(db, gameId)
	if err != nil {
		return err
	}

	if alreadySelected, err := selectSide(db, session, game, side); err != nil {
		return c.String(http.StatusInternalServerError, "Cannot select side")
	} else if alreadySelected {
		SaveSession(session, c)
		return c.Render(http.StatusOK, "select-side", SideSelectorData{
			XSelected:       game.XUserId.Valid,
			OSelected:       game.OUserId.Valid,
			AlreadySelected: alreadySelected,
		})
	}

	return GetBoard(ctx)
}

func CheckSideAndGetBoard(ctx Context) error {
		c := ctx.EchoCtx
		db := ctx.Db
		sess := GetSession(c)
		SaveSession(sess, c)
		gameId, err := strconv.Atoi(c.Param("gameId"))
		if err != nil {
			c.Logger().Error("could not parse gameId param", c.Param("gameId"))
			return c.String(http.StatusInternalServerError, "Cannot parse game id")
		}
		sess.Values["gameId"] = gameId
		game, err := repository.GetGame(db, gameId)
		if err != nil {
			c.Logger().Error("could not get game with id", gameId)
			return c.String(http.StatusInternalServerError, "Cannot find game")
		}
		sideId := fmt.Sprintf("side-%d", gameId)
		userId := sess.Values["userId"].(int)
		if int(game.XUserId.Int32) == userId {
			sess.Values[sideId] = "x"
		}
		if int(game.OUserId.Int32) == userId {
			sess.Values[sideId] = "o"
		}
		if sess.Values[sideId] == nil || sess.Values[sideId] == "spectator" {
			SaveSession(sess, c)
			return c.Render(http.StatusOK, "select-side", SideSelectorData{
				XSelected: game.XUserId.Valid,
				OSelected: game.OUserId.Valid,
				AlreadySelected: false,
			})
		}
		return GetBoard(ctx)
}

func GetBoard(ctx Context) error {
	c := ctx.EchoCtx
	db := ctx.Db
	session := GetSession(c)
	gameId := session.Values["gameId"].(int)
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

func selectSide(db *sqlx.DB, session *sessions.Session, game repository.Game, side string) (bool, error) {
	var err error
	userId := session.Values["userId"].(int)
	xSelected := game.XUserId.Valid
	if side == "x" && xSelected {
		return true, nil
	}
	oSelected := game.OUserId.Valid
	if side == "o" && oSelected {
		return true, nil
	}
	if side == "x" {
		err = repository.UpdateGameXUserId(db, game.Id, userId)
	} else if side == "o" {
		err = repository.UpdateGameOUserId(db, game.Id, userId)
	}
	if err != nil {
		return false, err
	}
	session.Values[fmt.Sprintf("side-%d", game.Id)] = side
	return false, nil
}
