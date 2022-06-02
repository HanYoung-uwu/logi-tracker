package database

import (
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"hanyoung/logi-tracker/pkg/utility"
)

type DataBaseManager struct {
	db *sql.DB
}

type Account struct {
	Name       string
	Permission int
	Id         int64
	Clan       string
}

var lock = &sync.Mutex{}
var singleton *DataBaseManager

func initDatabase() *sql.DB {
	m_db, _ := sql.Open("sqlite3",
		"test.sqlite3")
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS location (location TEXT PRIMARY KEY, time DATETIME, clan TEXT, code TEXT);
	CREATE TABLE IF NOT EXISTS account (name TEXT PRIMARY KEY, password TEXT, permission INTEGER, clan TEXT);
	CREATE TABLE IF NOT EXISTS salts (name TEXT PRIMARY KEY, salt TEXT);
	CREATE TABLE IF NOT EXISTS item (type TEXT, location TEXT, sum INTEGER, clan TEXT);
	`
	m_db.Exec(sqlStmt)
	return m_db
}

func GetInstance() *DataBaseManager {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			singleton = &DataBaseManager{initDatabase()}
		}
	}
	return singleton
}

func (manager *DataBaseManager) InsertItem(location string, item_type string, sum int, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare("insert into item(location, type, sum, clan) values(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(location, item_type, sum, clan)
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func (manager *DataBaseManager) CreateStockpile(location string, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare("insert into location(location, time, clan) values(?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(location, time.Now().Format(time.RFC3339), clan)
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func (manager *DataBaseManager) AddAccount(name string, password string, clan string, permission int) {
	// 0 is super admin, 1 is clan admin, 2 is ordinary memenber
	if permission < 0 || permission > 2 {
		log.Panic("unexpected permission value: ", permission)
	}

	salt := utility.RandBytes(256)

	hashedPassword, err := bcrypt.GenerateFromPassword(append(salt, password...), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	tx, err := manager.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare("insert into account(name, password, permission, clan) values(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, hashedPassword, permission, clan)
	if err != nil {
		log.Panic(err)
	}

	stmt, err = tx.Prepare("insert into salts(name, salt) values(?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, salt)
	if err != nil {
		log.Panic(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

var ErrorNoAccount = errors.New("no account found")
var ErrorIncorrectPassword = errors.New("incorrect password")

func (m *DataBaseManager) GetAndValidateAccount(name string, password string) (*Account, error) {
	stmt, err := m.db.Prepare("select id, password, permisson, clan from account where name = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	var id int64
	var hashedPassword string
	var clan string
	var permission int
	err = stmt.QueryRow(name).Scan(&id, &hashedPassword, &permission, &clan)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Panic(err)
		}
		return nil, ErrorNoAccount
	}

	stmt, err = m.db.Prepare("select salt from account where id = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	var salt string
	err = stmt.QueryRow(name).Scan(&salt)
	if err != nil {
		log.Panic(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), append([]byte(salt), password...))
	if err != nil {
		return nil, ErrorIncorrectPassword
	}
	return &Account{name, permission, id, clan}, nil
}
