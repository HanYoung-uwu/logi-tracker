package database

import (
	sql "database/sql"
	fmt "fmt"
	log "log"
	time "time"

	_ "github.com/mattn/go-sqlite3"
)

func Test() {
	db, _ := sql.Open("sqlite3",
		"test.sqlite3")
	defer db.Close()
	sqlStmt := `
	create table location (location TEXT NOT NULL PRIMARY KEY, time DATETIME);
	`
	db.Exec(sqlStmt)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into location(location, time) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(fmt.Sprintf("世界%03d", i), time.Now().Format(time.RFC3339))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
