package go_sqlite3_locked_tests

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"fmt"
)

type Sqlite struct {
	db *sql.DB
}

func (s *Sqlite) Init(name string) (err error) {
	path := fmt.Sprintf("./%s.db", name)
	os.Remove(path)
	s.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return
	}
	_, err = s.db.Exec( `
		CREATE TABLE IF NOT EXISTS test (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR NOT NULL DEFAULT ''
		)
	`)
	if err != nil {
		return
	}
	return
}

func (s *Sqlite) Insert(key string) (err error) {
	_, err = s.db.Exec(`
		INSERT INTO test(key) VALUES(?)
	`, key)
	return
}

func (s *Sqlite) Delete(key string) (err error) {
	_, err = s.db.Exec(`
		DELETE FROM test WHERE key=? 
	`, key)
	return
}

func (s *Sqlite) Update(id int, key string) (err error) {
	_, err = s.db.Exec(`
		UPDATE test SET key=? WHERE id=?
	`, key, id)
	return
}

func (s *Sqlite) Select(id int) (key string, err error) {
	res, err := s.db.Query(`
		  SELECT key FROM test WHERE id=?
	`, id)
	if err != nil {
		return
	}
	defer res.Close()
	if res.Next() {
		err = res.Scan(&key)
	}
	return
}

func (s *Sqlite) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}