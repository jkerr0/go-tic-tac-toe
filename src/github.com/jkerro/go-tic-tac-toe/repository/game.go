package repository

import "github.com/jmoiron/sqlx"

type Game struct {
	Id int `db:"id"`
	Name string `db:"name"`
}

func GetGames(db *sqlx.DB) ([]Game, error) {
	games := []Game{}
	err := db.Select(&games, "SELECT * FROM game ORDER BY name")
	return games, err
}

func InsertGame(db *sqlx.DB, name string) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO game(name) VALUES (?)", name)
	return tx.Commit()
}

func DeleteGame(db *sqlx.DB, id int) error {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM game WHERE id=?", id)
	tx.MustExec("DELETE FROM move WHERE game_id=?", id)
	return tx.Commit()
}