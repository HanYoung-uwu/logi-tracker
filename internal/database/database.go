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
	Clan       string
}

type StockpileItem struct {
	ItemType string `json:"item"`
	Location string `json:"location"`
	Size     int    `json:"size"`
}

type Location struct {
	Location string
	Time     time.Time
	Code     string
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
	CREATE TABLE IF NOT EXISTS item (type TEXT, location TEXT, size INTEGER, clan TEXT);
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

func (manager *DataBaseManager) InsertItem(location string, item_type string, size int, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare("insert into item(location, type, size, clan) values(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(location, item_type, size, clan)

	stmt, err = tx.Prepare("update location set time=? where location=? and clan=?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(time.Now().Format(time.RFC3339), location, clan)

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func (manager *DataBaseManager) CreateStockpile(location string, code string, clan string) {
	tx, err := manager.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare("insert into location(location, time, clan, code) values(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(location, time.Now().Format(time.RFC3339), clan, code)
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
	stmt, err := m.db.Prepare("select password, permission, clan from account where name = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	var hashedPassword string
	var clan string
	var permission int
	err = stmt.QueryRow(name).Scan(&hashedPassword, &permission, &clan)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Panic(err)
		}
		return nil, ErrorNoAccount
	}

	stmt, err = m.db.Prepare("select salt from salts where name = ?")
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
	return &Account{name, permission, clan}, nil
}

func (m *DataBaseManager) GetAllItems(account *Account) []StockpileItem {
	clan := account.Clan

	db := m.db

	stmt, err := db.Prepare("select type, location, size from item where clan = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(clan)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	resultArray := make([]StockpileItem, 0, 30)

	for rows.Next() {
		var itemType string
		var location string
		var size int
		rows.Scan(&itemType, &location, &size)
		resultArray = append(resultArray, StockpileItem{itemType, location, size})
	}
	return resultArray
}

func (m *DataBaseManager) GetAllLocations(clan string) []Location {
	db := m.db

	stmt, err := db.Prepare("select location, time, code from location where clan = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(clan)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	resultArray := make([]Location, 0, 30)

	for rows.Next() {
		var location string
		var time time.Time
		var code string
		rows.Scan(&location, &time, &code)
		resultArray = append(resultArray, Location{location, time, code})
	}
	return resultArray
}
