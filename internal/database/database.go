package database

import (
	"crypto/rand"
	"database/sql"
	"log"
	"math/big"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type DataBaseManager struct {
	db *sql.DB
}

var lock = &sync.Mutex{}
var db *DataBaseManager

func initDatabase() *sql.DB {
	m_db, _ := sql.Open("sqlite3",
		"test.sqlite3")
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS location (id INTEGER PRIMARY KEY AUTOINCREMENT, location TEXT, time DATETIME, clan TEXT);
	CREATE TABLE IF NOT EXISTS account (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, password TEXT, admin INTEGER, clan TEXT);
	CREATE TABLE IF NOT EXISTS salts (id INTEGER PRIMARY KEY, salt TEXT);
	CREATE TABLE IF NOT EXISTS item (id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT, location TEXT, sum INTEGER, clan TEXT);
	`
	m_db.Exec(sqlStmt)
	return m_db
}

func GetInstance() *DataBaseManager {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			db = &DataBaseManager{initDatabase()}
		}
	}
	return db
}

func (manager DataBaseManager) InsertItem(location string, item_type string, sum int, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into item(location, type, sum, clan) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(location, item_type, sum, clan)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (manager DataBaseManager) CreateStockpile(location string, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into location(location, time, clan) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(location, time.Now().Format(time.RFC3339), clan)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (manager DataBaseManager) AddAccount(name string, password string, clan string, admin int) {
	// 0 is super admin, 1 is clan admin, 2 is ordinary memenber
	if admin < 0 || admin > 2 {
		log.Panic("unexpected admin value: ", admin)
	}

	salt := randBytes(50)

	hashedPassword, err := bcrypt.GenerateFromPassword(append(salt, password...), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := manager.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into account(name, password, admin, clan) values(?, ?, ?, ?) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, hashedPassword, admin, clan)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	tx, err = manager.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err = tx.Prepare("insert into salts(id, salt) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, salt)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func randBytes(l int) []byte {
	store := "asdfghjkl;'qwertyuiop[]1234567890-=zxcvbnm,./QWERTYUIOP{}ASDFGHJKL:\\\"|ZXCVBNM<>?!@#$%^&*()_+"
	maxLen := len(store)
	result := make([]byte, l)
	i := 0
	for {
		if i == l {
			break
		}
		p, err := rand.Int(rand.Reader, big.NewInt(int64(maxLen)))
		if err != nil {
			log.Fatal(err)
		}
		result[i] = store[p.Int64()]
		i += 1
	}
	return result
}
