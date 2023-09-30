package repository

import "github.com/jmoiron/sqlx"

type User struct {
	Id int `db:"id"`
}

func CreateUser(db *sqlx.DB) (int, error) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO app_user DEFAULT VALUES")
	var id int
	tx.QueryRow("SELECT max(id) FROM app_user").Scan(&id)
	if err := tx.Commit(); err != nil {
		return -1, err
	}
	return id, nil 
}
