package repository

import "github.com/jmoiron/sqlx"

func CreateDatabase(db *sqlx.DB) {
	tx := db.MustBegin()
	createTableGame(tx)
	createTableMove(tx)
	createTableUser(tx)
	tx.Commit()
}

func createTableGame(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS game (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			x_user_id INTEGER NULL,
			o_user_id INTEGER NULL,
			FOREIGN KEY (x_user_id) REFERENCES user(id),
			FOREIGN KEY (o_user_id) REFERENCES user(id)
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

func createTableUser(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT
		)
	`)
}
