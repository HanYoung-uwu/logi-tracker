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

// 0 is super admin, 1 is clan admin, 2 is ordinary memenber, 3 is temporary account for invitation links
type Account struct {
	Name       string
	Permission int
	Clan       string
}

var NormalAccount = 2
var AdminAccount = 0
var ClanAdminAccount = 1
var InvitationLinkAccount = 3

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
		utility.DatabasePath)
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

var ErrorUnableToUpdateItem = errors.New("can't retrieve item, not enough in stockpile to retrieve")

// a negative size means retrieval
func (m *DataBaseManager) InsertOrUpdateItem(location string, item string, size int, clan string) error {
	tx, err := m.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	// check if we already have this item
	stmt, err := tx.Prepare("select size from item where clan = ? and location = ? and type = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	var number int
	err = stmt.QueryRow(clan, location, item).Scan(&number)
	if err != nil {
		// new item
		if size < 0 {
			return ErrorUnableToUpdateItem
		}
		stmt, err := tx.Prepare("insert into item(location, type, size, clan) values(?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(location, item, size, clan)

		// update the stockpile's time as well
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
		return nil
	} else {
		updatedSize := number + size
		if updatedSize > 0 {
			stmt, err := m.db.Prepare("update item set size = ? where clan = ? and location = ? and type = ?")
			if err != nil {
				log.Panic(err)
			}
			defer stmt.Close()
			stmt.Exec(updatedSize, clan, location, item)
			return nil
		} else if updatedSize == 0 {
			stmt, err := m.db.Prepare("delete from item where clan = ? and location = ? and type = ?")
			if err != nil {
				log.Panic(err)
			}
			defer stmt.Close()
			stmt.Exec(updatedSize, clan, location, item)
			return nil
		} else {
			return ErrorUnableToUpdateItem
		}
	}
}

func (m *DataBaseManager) CreateStockpile(location string, code string, clan string) {
	tx, err := m.db.Begin()
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

func (m *DataBaseManager) AddAccount(name string, password string, clan string, permission int) {
	// 0 is super admin, 1 is clan admin, 2 is ordinary memenber
	if permission < 0 || permission > 2 {
		log.Panic("unexpected permission value: ", permission)
	}

	salt := utility.RandBytes(256)

	hashedPassword, err := bcrypt.GenerateFromPassword(append(salt, password...), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	tx, err := m.db.Begin()
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
