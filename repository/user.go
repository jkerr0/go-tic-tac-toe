package repository

import "github.com/jmoiron/sqlx"

type User struct {
	Id int `db:"id"`
}

func CreateUser(db *sqlx.DB) (int, error) {
	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO user DEFAULT VALUES")
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	err = tx.Commit()
	return int(id), err
}
