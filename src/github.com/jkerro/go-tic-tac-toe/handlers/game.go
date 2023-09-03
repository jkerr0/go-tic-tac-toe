package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
)

type move struct {
	Row int
	Col int
}
type moveResponse struct {
	Row          int
	Col          int
	Info         string
	BoardElement logic.BoardElement
}

func HandleMoveMessage(ctx Context, message []byte, gameId int, channel *Channel) {
	c := ctx.EchoCtx

	parsedMove, err := parseMessage(message)
	if err != nil {
		c.Logger().Error("could not parse message", err)
		return
	}
	col := parsedMove.Col
	row := parsedMove.Row

	move := &repository.Move{
		Col:    col,
		Row:    row,
		GameId: gameId}
	err = move.Insert(ctx.Db)
	if err != nil {
		c.Logger().Error("could not insert move", err)
		return
	}

	response, err := getResponse(ctx, move, gameId)
	if err != nil {
		c.Logger().Error("could not generate response", err)
		return
	}
	channel.Broadcast(response)
}

func parseMessage(message []byte) (move, error) {
	var data map[string]interface{}
	err := json.Unmarshal(message, &data)
	empty := move{}
	if err != nil {
		return empty, err
	}
	col, err := strconv.Atoi(data["col"].(string))
	if err != nil {
		return empty, err
	}
	row, err := strconv.Atoi(data["row"].(string))
	if err != nil {
		return empty, err
	}
	return move{
		Row: row,
		Col: col,
	}, nil
}

func checkGameWin(db *sqlx.DB, gameId int) (bool, error) {
	moves, err := repository.GetMoves(db, gameId)
	if err != nil {
		return false, err
	}
	return logic.CheckWin(logic.GetBoard(moves)), nil
}

func getResponse(ctx Context, move *repository.Move, gameId int) (string, error) {
	c := ctx.EchoCtx
	turn := logic.GetTurnElement(move.Inx)
	win, err := checkGameWin(ctx.Db, gameId)
	if err != nil {
		return "", err
	}
	var info string
	if win {
		info = fmt.Sprintf("%s wins", turn)
	}

	var responseBuf bytes.Buffer
	responseData := moveResponse{
		Col:          move.Col,
		Row:          move.Row,
		Info:         info,
		BoardElement: turn,
	}
	c.Echo().Renderer.Render(&responseBuf, "move-response", responseData, c)
	return responseBuf.String(), nil
}
