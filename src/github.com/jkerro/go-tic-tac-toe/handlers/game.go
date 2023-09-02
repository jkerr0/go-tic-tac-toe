package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jkerro/go-tic-tac-toe/logic"
	"github.com/jkerro/go-tic-tac-toe/repository"
)

func HandleMoveMessage(ctx Context, message []byte, gameId int, channel *Channel) {
	c := ctx.EchoCtx
	var data map[string]interface{}

	err := json.Unmarshal(message, &data)

	if err != nil {
		c.Logger().Error("could not unmarshal json: %s\n", err)
		return
	}
	col, err := strconv.Atoi(data["col"].(string))
	if err != nil {
		c.Logger().Error("could not parse column index")
		return
	}
	row, err := strconv.Atoi(data["row"].(string))
	if err != nil {
		c.Logger().Error("could not parse row index")
		return
	}

	move := &repository.Move{
		Col:    col,
		Row:    row,
		GameId: gameId}
	err = move.Insert(ctx.Db)
	if err != nil {
		c.Logger().Error(err)
		return
	}

	turn := logic.GetTurnElement(move.Inx)
	var info string
	moves, err := repository.GetMoves(ctx.Db, gameId)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	win := logic.CheckWin(logic.GetBoard(moves))
	if win {
		info = fmt.Sprintf("%s wins", turn)
	}

	channel.Broadcast(fmt.Sprintf("<button id=\"row-%d-col-%d\">%s</button><div id=\"info\">%s</div>", row, col, turn, info))
}
