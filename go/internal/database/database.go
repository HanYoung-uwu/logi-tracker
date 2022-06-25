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

var FactionW = 1
var FactionC = 2

type DataBaseManager struct {
	db              *sql.DB
	itemLock        *sync.Map
	locationLock    *sync.Map
	registeredNames *sync.Map
	registeredClans *sync.Map
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

type HistoryRecord struct {
	// 0 add, 1 retrieve, 2 delete, 3 set,
	// 4 add stockpile, 5 delete stockpile
	Action   int
	Time     time.Time
	ItemType string
	Location string
	Size     int
	Clan     string
	User     string
}

type Token struct {
	Value      string
	ExpireTime time.Time
	Account    *Account
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
	CREATE TABLE IF NOT EXISTS history (action INTEGER, user TEXT, clan TEXT, type TEXT, size INTEGER, location TEXT, time DATETIME);
	CREATE INDEX IF NOT EXISTS idx_item_history on history (clan, location);
	CREATE TABLE IF NOT EXISTS tokens (token TEXT, expire_time DATETIME, account_name TEXT);
	CREATE TABLE IF NOT EXISTS clan (name TEXT PRIMARY KEY, faction INTEGER);
	`
	m_db.Exec(sqlStmt)
	return m_db
}

func GetInstance() *DataBaseManager {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			singleton = &DataBaseManager{initDatabase(),
				&sync.Map{},
				&sync.Map{},
				&sync.Map{},
				&sync.Map{}}
			singleton.loadInfoToMemory()
		}
	}
	return singleton
}

func (m *DataBaseManager) loadInfoToMemory() {
	stmt, err := m.db.Prepare("select name from account")
	if err != nil {
		log.Panic(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Panic(err)
		}
		m.registeredNames.Store(name, true)
	}
	stmt, err = m.db.Prepare("select name, faction from clan")
	if err != nil {
		log.Panic(err)
	}
	rows, err = stmt.Query()
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var name string
		var faction int
		err = rows.Scan(&name, &faction)
		if err != nil {
			log.Panic(err)
		}
		m.registeredClans.Store(name, faction)
	}
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
func (m *DataBaseManager) InsertOrUpdateItem(location string, item string, size int, clan string, user string) error {
	err := m._InsertOrUpdateItem(location, item, size, clan)
	if err != nil {
		return err
	} else {
		go func(location string, item string, size int, clan string) {
			// update the stockpile's time
			stmt, err := m.db.Prepare("update location set time=? where location=? and clan=?")
			if err != nil {
				log.Panic(err)
			}
			defer stmt.Close()
			stmt.Exec(time.Now().Format(time.RFC3339), location, clan)

			// log to history
			stmt, err = m.db.Prepare("insert into history(action, user, clan, type, size, location, time) values(?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Panic(err)
			}
			defer stmt.Close()
			var action int
			if size > 0 {
				action = 0
			} else {
				action = 1
			}
			_, err = stmt.Exec(action, user, clan, item, size, location, time.Now().Format(time.RFC3339))
			if err != nil {
				log.Panic(err)
			}
		}(location, item, size, clan)
	}
	return nil
}
func (m *DataBaseManager) _InsertOrUpdateItem(location string, item string, size int, clan string) error {
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

		stmt, err := m.db.Prepare("insert into item(location, type, size, clan) values(?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(location, item, size, clan)

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

func (m *DataBaseManager) DeleteItem(location string, item string, clan string, user string) {
	lock := m.getItemLock(location, item, clan)
	lock.Lock()
	defer lock.Unlock()
	stmt, err := m.db.Prepare("delete from item where clan = ? and location = ? and type = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(clan, location, item)

	if err != nil {
		log.Panic(err)
	}
	go func(location string, item string, clan string) {
		// update the stockpile's time
		stmt, err := m.db.Prepare("update location set time=? where location=? and clan=?")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(time.Now().Format(time.RFC3339), location, clan)

		// log to history
		stmt, err = m.db.Prepare("insert into history(action, user, clan, type, location, time) values(?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		action := 2
		_, err = stmt.Exec(action, user, clan, item, location, time.Now().Format(time.RFC3339))
		if err != nil {
			log.Panic(err)
		}
	}(location, item, clan)
}

func (m *DataBaseManager) SetItem(location string, item string, size int, clan string, user string) {
	lock := m.getItemLock(location, item, clan)
	lock.Lock()
	defer lock.Unlock()
	stmt, err := m.db.Prepare("update item set size = ? where clan = ? and location = ? and type = ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(size, clan, location, item)
	go func(location string, item string, clan string) {
		// update the stockpile's time
		stmt, err := m.db.Prepare("update location set time=? where location=? and clan=?")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(time.Now().Format(time.RFC3339), location, clan)

		// log to history
		stmt, err = m.db.Prepare("insert into history(action, user, clan, type, size, location, time) values(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		action := 3
		_, err = stmt.Exec(action, user, clan, item, size, location, time.Now().Format(time.RFC3339))
		if err != nil {
			log.Panic(err)
		}
	}(location, item, clan)
}

func (m *DataBaseManager) CreateStockpile(location string, code string, clan string, user string) {
	stmt, err := m.db.Prepare("insert into location(location, time, clan, code) values(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(location, time.Now().Format(time.RFC3339), clan, code)
	if err != nil {
		log.Panic(err)
	}

	go func(location string, clan string, user string) {
		stmt, err := m.db.Prepare("insert into history(action, location, time, clan, user) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(4, location, time.Now().Format(time.RFC3339), clan, user)
		if err != nil {
			log.Panic(err)
		}
	}(location, clan, user)
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
	m.registeredNames.Store(name, true)
	if permission == ClanAdminAccount {
		m.registeredClans.Store(clan, true)
		_, err = m.db.Exec("insert into clan(name, faction) values(?, 0)", clan)
		if err != nil {
			log.Print(err)
		}
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

func (m *DataBaseManager) DeleteStockpile(location string, clan string, user string) error {
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
	stmt, err = m.db.Prepare("delete from item where location = ?")
	if err != nil {
		log.Panic(err)
	}

	_, err = stmt.Exec(location)
	if err != nil {
		log.Panic(err)
	}
	go func(location string, clan string, user string) {
		stmt, err := m.db.Prepare("insert into history(action, location, time, clan, user) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Panic(err)
		}
		defer stmt.Close()
		stmt.Exec(5, location, time.Now().Format(time.RFC3339), clan, user)
		if err != nil {
			log.Panic(err)
		}
	}(location, clan, user)
	return nil
}

func (m *DataBaseManager) GetClanHistory(clan string, limit ...int) []HistoryRecord {
	var queryLimit int
	if len(limit) == 0 {
		queryLimit = 30
	} else {
		queryLimit = limit[0]
	}

	stmt, err := m.db.Prepare("select action, user, type, size, location, time from history where clan = ? order by time desc limit ?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(clan, queryLimit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	resultArray := make([]HistoryRecord, 0, queryLimit)

	for rows.Next() {
		var action int
		var user string
		var item sql.NullString
		var size sql.NullInt32
		var time time.Time
		var location string
		rows.Scan(&action, &user, &item, &size, &location, &time)
		resultArray = append(resultArray, HistoryRecord{
			Action:   action,
			User:     user,
			ItemType: item.String,
			Size:     int(size.Int32),
			Location: location,
			Time:     time,
		})
	}
	return resultArray
}

func (m *DataBaseManager) RefreshStockpile(location string, clan string) {
	stmt, err := m.db.Prepare("update location set time=? where location=? and clan=?")
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format(time.RFC3339), location, clan)
	if err != nil {
		log.Panic(err)
	}
}

func (m *DataBaseManager) IsNameExist(name string) bool {
	_, found := m.registeredNames.Load(name)
	return found
}

func (m *DataBaseManager) IsClanExist(name string) bool {
	_, found := m.registeredClans.Load(name)
	return found
}

func (m *DataBaseManager) SaveTokens(tokens []interface{}) {
	stmt, err := m.db.Prepare("insert into tokens(token, expire_time, account_name) values(?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}
	for _, token := range tokens {
		_token, ok := token.(Token)
		if !ok {
			log.Panic("passed non Token array!")
		}
		_, err = stmt.Exec(_token.Value, _token.ExpireTime, _token.Account.Name)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (m *DataBaseManager) DeleteTokens(tokens []interface{}) {
	stmt, err := m.db.Prepare("delete from tokens where token = ?")
	if err != nil {
		log.Panic(err)
	}
	for _, token := range tokens {
		_token, ok := token.(string)
		if !ok {
			log.Panic("passed non string array!")
		}
		_, err = stmt.Exec(_token)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (m *DataBaseManager) LoadTokens() []Token {
	stmt := `
	select t.token, t.expire_time, t.account_name, a.clan, a.permission
	from tokens as t
	join account as a
	where a.name = t.account_name
	`
	rows, err := m.db.Query(stmt)
	if err != nil {
		log.Panic(err)
	}
	result := make([]Token, 0, 50)
	for rows.Next() {
		var token string
		var expire_time time.Time
		var account_name string
		var clan string
		var permission int
		rows.Scan(&token, &expire_time, &account_name, &clan, &permission)
		account := &Account{account_name, permission, clan}
		result = append(result, Token{token, expire_time, account})
	}
	return result
}

func (m *DataBaseManager) GetClanMembers(clan string) []Account {
	result := make([]Account, 0, 20)
	row, err := m.db.Query("select name, permission from account where clan = ?", clan)
	if err != nil {
		log.Panic(err)
	}
	for row.Next() {
		var name string
		var permission int
		err = row.Scan(&name, &permission)
		if err != nil {
			log.Panic(err)
		}
		result = append(result, Account{name, permission, clan})
	}
	return result
}

func (m *DataBaseManager) PromoteClanMember(clan string, name string) {
	_, err := m.db.Exec("update account set permission = 1 where name = ? and clan = ?", name, clan)
	if err != nil {
		log.Panic(err)
	}
}

func (m *DataBaseManager) KickClanMember(clan string, name string) {
	_, err := m.db.Exec("delete from account where name = ? and clan = ?", name, clan)
	if err != nil {
		log.Panic(err)
	}
}
