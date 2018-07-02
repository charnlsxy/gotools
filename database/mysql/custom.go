package mysql

import "database/sql"

var testDb *sql.DB

func GetOneUsableDb() *sql.DB {
	return testDb
}