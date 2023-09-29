package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Game struct {
	Id      int           `db:"id"`
	Name    string        `db:"name"`
	XUserId sql.NullInt32 `db:"x_user_id"`
	OUserId sql.NullInt32 `db:"o_user_id"`
}

func GetGames(db *sqlx.DB) ([]Game, error) {
	games := []Game{}
	err := db.Select(&games, "SELECT * FROM game ORDER BY name")
	return games, err
}

func GetGame(db *sqlx.DB, gameId int) (Game, error) {
	games := []Game{}
	err := db.Select(&games, "SELECT * FROM game WHERE id=$1", gameId)
	return games[0], err
}

func UpdateGameXUserId(db *sqlx.DB, gameId int, xUserId int) error {
	tx := db.MustBegin()
	tx.MustExec(`
		UPDATE game
		SET x_user_id=$1
		WHERE id=$2`, xUserId, gameId)
	return tx.Commit()
}

func UpdateGameOUserId(db *sqlx.DB, gameId int, oUserId int) error {
	tx := db.MustBegin()
	tx.MustExec(`
		UPDATE game
		SET o_user_id=$1
		WHERE id=$2`, oUserId, gameId)
	return tx.Commit()
}

func InsertGame(db *sqlx.DB, name string) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO game (name) VALUES ($1)", name)
	return tx.Commit()
}

func DeleteGame(db *sqlx.DB, id int) error {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM game WHERE id=$1", id)
	tx.MustExec("DELETE FROM move WHERE game_id=$1;", id)
	return tx.Commit()
}
