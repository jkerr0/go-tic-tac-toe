package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Move struct {
	GameId int `db:"game_id"`
	Col    int `db:"col"`
	Row    int `db:"row"`
	Inx    int `db:"inx"`
}

func (m *Move) Insert(db *sqlx.DB) error {
	tx := db.MustBegin()
	maxInx, err := getMaxIndex(db, m.GameId)
	if err != nil {
		return err
	}
	m.Inx = maxInx + 1
	tx.NamedExec("INSERT INTO move (game_id, col, row, inx) VALUES (:game_id, :col, :row, :inx)", m)
	return tx.Commit()
}

func getMaxIndex(db *sqlx.DB, gameId int) (int, error) {
	index := make([]sql.NullInt32, 1)
	err := db.Select(&index, "SELECT max(inx) FROM move WHERE game_id=?", gameId)
	if len(index) == 0 || !index[0].Valid {
		return 0, err
	}
	return int(index[0].Int32), err
}

func GetMoves(db *sqlx.DB, gameId int) ([]Move, error) {
	moves := []Move{}
	err := db.Select(&moves, "SELECT * FROM move WHERE game_id=?", gameId)
	return moves, err
}
