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
	db       *sql.DB
	itemLock *sync.Map
	locationLock *sync.Map
}

// 0 is super admin, 1 is clan admin, 2 is ordinary memenber, 3 is temporary account for invitation links,
// 4 is clan admin invitation links

type Account struct {
	Name       string
	Permission int
	Clan       string
}

var NormalAccount = 2
var AdminAccount = 0
var ClanAdminAccount = 1
var InvitationLinkAccount = 3
var ClanAdminInvitationLinkAccount = 4

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
			singleton = &DataBaseManager{initDatabase(), &sync.Map{}, &sync.Map{}}
		}
	}
	return singleton
}

func (m *DataBaseManager) getItemLock(location string, item string, clan string) *sync.Mutex {
	key := item + "###" + location + "###" + clan
	val, _ := m.itemLock.LoadOrStore(key, &sync.Mutex{})
	lock, ok := val.(*sync.Mutex)
	if !ok {
		log.Panic("val is not a lock")
	}
	return lock
}

func (m *DataBaseManager) getLocationLock(location string, clan string) *sync.Mutex {
	key := location + "###" + clan
	val, _ := m.locationLock.LoadOrStore(key, &sync.Mutex{})
	lock, ok := val.(*sync.Mutex)
	if !ok {
		log.Panic("val is not a lock")
	}
	return lock
}

var ErrorUnableToUpdateItem = errors.New("can't retrieve item, not enough in stockpile to retrieve")

// a negative size means retrieval
func (m *DataBaseManager) InsertOrUpdateItem(location string, item string, size int, clan string) error {
	lock := m.getItemLock(location, item, clan)
	lock.Lock()
	defer lock.Unlock()

	// check if we already have this item
	stmt, err := m.db.Prepare("select size from item where clan = ? and location = ? and type = ?")
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

		tx, err := m.db.Begin()
		if err != nil {
			log.Panic(err)
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

func (m *DataBaseManager) AddAccount(name string, password string, clan string, permission int) error {
	// 0 is super admin, 1 is clan admin, 2 is ordinary memenber
	if permission < 0 || permission > 2 {
		log.Panic("unexpected permission value: ", permission)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
		log.Println(err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	return nil
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

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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

var ErrorLocationNotExists = errors.New("location doesn't exists")
func (m *DataBaseManager) DeleteStockpile(location string, clan string) error {
	lock := m.getLocationLock(location, clan)
	lock.Lock()
	defer lock.Unlock()
	stmt, err := m.db.Prepare("delete from location where location.location = ?")
	if err != nil {
		log.Panic(err)
	}

	result, err := stmt.Exec(location)
	if err != nil {
		log.Panic(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Panic(err)
	}
	if affected == 0 {
		return ErrorLocationNotExists
	}
	return nil
}