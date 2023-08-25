package main

import (
	"log"

	"github.com/jkerro/go-tic-tac-toe/repository"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)
func main() {
	db, err := sqlx.Connect("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	log.Println("Database connected")
	repository.InsertGame(db, "newgame")
	games, err := repository.GetGames(db)
	for _, game := range games {
		log.Println(game.Name)
	}
}