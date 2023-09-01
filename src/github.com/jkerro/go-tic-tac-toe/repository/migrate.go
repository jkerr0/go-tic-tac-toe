package repository

import "github.com/jmoiron/sqlx"

func CreateDatabase(db *sqlx.DB) {
	tx := db.MustBegin()
	createTableGame(tx)
	createTableMove(tx)
	tx.Commit()
}

func createTableGame(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS game (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)
	`)
}

func createTableMove(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS move (
			game_id INTEGER,
			col INTEGER,
			row INTEGER,
			inx INTEGER,
			PRIMARY KEY (game_id, col, row, inx),
			FOREIGN KEY(game_id) REFERENCES game(id)
		)
	`)
}