package store

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	myDb *sql.DB
	once sync.Once
)

func getDb() *sql.DB {
	once.Do(func() {
		dbAddr := "testuser:123456@tcp(127.0.0.1:3306)/test"
		db, err := sql.Open("mysql", dbAddr)
		checkErr(err)
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(10)
		pingErr := db.Ping()
		if pingErr != nil {
			panic(pingErr)
		}
		myDb = db
	})
	return myDb
}

func GetDb() *sql.DB {
	return getDb()
}

// panic and ends the program here, when can't connect to base db
func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// GetFromResult get generated id and rows affected and possible error
func GetFromResult(result sql.Result, passError error) (int64, int64, error) {
	if passError != nil {
		return 0, 0, errors.New("result should not be null")
	}
	if result == nil {
		return 0, 0, errors.New("result should not be null")
	}
	var finalError error = nil
	id, err := result.LastInsertId()
	if err != nil {
		finalError = err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		finalError = err
	}
	return id, affected, finalError
}

// nullable data insertion

// NewNullString nullable string
func NewNullString(i interface{}) sql.NullString {
	if i == nil {
		return sql.NullString{}
	}
	if _, ok := i.(*string); !ok {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *i.(*string),
		Valid:  true,
	}
}

// NewNullInt nullable int
func NewNullInt(i interface{}) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	if _, ok := i.(*int64); !ok {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *i.(*int64),
		Valid: true,
	}
}

// NewNullFloat nullable float
func NewNullFloat(i interface{}) sql.NullFloat64 {
	if i == nil {
		return sql.NullFloat64{}
	}
	if _, ok := i.(*float64); !ok {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *i.(*float64),
		Valid:   true,
	}
}

// NewNullBool nullable float
func NewNullBool(i interface{}) sql.NullBool {
	if i == nil {
		return sql.NullBool{}
	}
	if _, ok := i.(*bool); !ok {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *i.(*bool),
		Valid: true,
	}
}
