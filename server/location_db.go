package server

import (
	"github.com/SUCHMOKUO/falcon-ws/util"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"path/filepath"
)

var (
	db *sql.DB
	insert *sql.Stmt
)

func init() {
	prepareDB()
	initTable()
	prepareStmt()
}

func prepareDB() {
	path := filepath.Join(util.GetCurrentPath(), "server.db")
	var err error

	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalln("open db error:", err)
	}
}

func initTable() {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS locations(
			host TEXT PRIMARY KEY,
			location TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalln("prepare table error:", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln("create table error:", err)
	}
}

func prepareStmt() {
	var err error
	insert, err = db.Prepare(`INSERT INTO locations(host, location) VALUES(?, ?)`)
	if err != nil {
		log.Fatalln("prepare insert error:", err)
	}
}

func getLocationFromDB(host string) (string, bool) {
	location := new(string)
	err := db.QueryRow("SELECT location FROM locations WHERE host = ?", host).Scan(location)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		log.Println("query error:", err)
		return "", false
	}
	return *location, true
}

func setLocationToDB(host, location string) {
	_, err := insert.Exec(host, location)
	if err != nil {
		log.Println("insert error:", err)
	}
}
