package repository

import "github.com/jmoiron/sqlx"

func CreateDatabase(db *sqlx.DB) {
	tx := db.MustBegin()
	createTableUser(tx)
	createTableGame(tx)
	createTableMove(tx)
	tx.Commit()
}

func createTableGame(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS game (
			id int4 GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name varchar(200) NOT NULL,
			x_user_id int4,
			o_user_id int4,
			FOREIGN KEY (x_user_id) REFERENCES app_user(id),
			FOREIGN KEY (o_user_id) REFERENCES app_user(id)
		)
	`)
}

func createTableMove(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS move (
			game_id int4,
			col int4,
			row int4,
			inx int4,
			PRIMARY KEY (game_id, col, row, inx),
			FOREIGN KEY(game_id) REFERENCES game(id)
		)
	`)
}

func createTableUser(tx *sqlx.Tx) {
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS app_user (
			id int4 GENERATED ALWAYS AS IDENTITY PRIMARY KEY
		)
	`)
}
