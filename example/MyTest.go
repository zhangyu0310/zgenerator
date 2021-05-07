// Package handler is generate by zgenerator (Author:poppinzhang).
package handler

import "database/sql"
import "fmt"
import "github.com/jmoiron/sqlx"

// TestTable from Database testdb, Table TestTable
type TestTable struct {
	ID    int64          `db:"ID"`    // UserID
	Name  sql.NullString `db:"Name"`  // User Name
	Money sql.NullInt64  `db:"Money"` // Money
}

// SelectTestTable select SQL for TestTable
var SelectTestTable = "SELECT * FROM %v WHERE ID=%v"

// InsertTestTable insert SQL for TestTable
var InsertTestTable = "INSERT INTO %s (Name, Money) VALUES (:Name, :Money) ON DUPLICATE KEY UPDATE Name=:Name, Money=:Money"

// DeleteTestTable delete SQL for TestTable
var DeleteTestTable = "DELETE FROM %v WHERE ID=%v"

// TableName table name from user input.
var TableName = "TestTable"

func Select(db *sqlx.DB, dest *[]TestTable, keys ...interface{}) error {
	return SelectWithTableName(db, dest, TableName, keys)
}

func SelectWithTableName(db *sqlx.DB, dest *[]TestTable, tableName string, keys ...interface{}) error {
	var tmpKeys []interface{}
	tmpKeys = append(tmpKeys, tableName)
	tmpKeys = append(tmpKeys, keys...)
	sqlStr := fmt.Sprintf(SelectTestTable, tmpKeys...)
	return db.Select(dest, sqlStr)
}

func Insert(db *sqlx.DB, src TestTable) (sql.Result, error) {
	return InsertWithTableName(db, TableName, src)
}

func InsertWithTableName(db *sqlx.DB, tableName string, src TestTable) (sql.Result, error) {
	sqlStr := fmt.Sprintf(InsertTestTable, tableName)
	return db.NamedExec(sqlStr, src)
}

func Delete(db *sqlx.DB, keys ...interface{}) (sql.Result, error) {
	return DeleteWithTableName(db, TableName, keys)
}

func DeleteWithTableName(db *sqlx.DB, tableName string, keys ...interface{}) (sql.Result, error) {
	var tmpKeys []interface{}
	tmpKeys = append(tmpKeys, tableName)
	tmpKeys = append(tmpKeys, keys...)
	sqlStr := fmt.Sprintf(DeleteTestTable, tmpKeys...)
	return db.Exec(sqlStr)
}
