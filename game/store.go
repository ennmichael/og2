package game

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func Init() (Store, error) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		return Store{}, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS games (
		user TEXT PRIMARY KEY,
		game_json TEXT NOT NULL
	)`)
	if err != nil {
		return Store{}, err
	}
	return Store{db}, nil
}

// These methods probably shouldn't panic, and should propagate the errors
// instead, and either handled or displayed more elegantly to the user.

func (s Store) SaveGame(state State) {
	j, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	_, err = s.db.Exec(
		`INSERT INTO games (user, game_json) VALUES (?1, ?2)
			ON CONFLICT (user) DO UPDATE SET game_json = ?2`,
		string(state.User),
		string(j))
	if err != nil {
		panic(err)
	}
}

func (s Store) LoadGames() []State {
	games := []State{}
	rows, err := s.db.Query("SELECT (user, game_json) FROM games")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var user string
		var gameJSON string
		rows.Scan(&user, &gameJSON)
		var state State
		err := json.Unmarshal([]byte(gameJSON), &state)
		if err != nil {
			panic(err)
		}
		games = append(games, state)
	}
	return games
}
