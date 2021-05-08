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

// SelectTestTableSQL select SQL for TestTable
var SelectTestTableSQL = "SELECT * FROM `%v` WHERE `ID`='%v'"

// InsertTestTableSQL insert SQL for TestTable
var InsertTestTableSQL = "INSERT INTO `%v` (`Name`, `Money`) VALUES (:Name, :Money) ON DUPLICATE KEY UPDATE `Name`=:Name, `Money`=:Money"

// DeleteTestTableSQL delete SQL for TestTable
var DeleteTestTableSQL = "DELETE FROM `%v` WHERE `ID`='%v'"

// TestTableTableName table name from user input.
var TestTableTableName = "TestTable"

// SelectTestTable use SelectTestTableSQL select table TestTable
func SelectTestTable(db *sqlx.DB, dest *[]TestTable, keys ...interface{}) error {
	return SelectTestTableWithTableName(db, dest, TestTableTableName, keys)
}

// SelectTestTableWithTableName use SelectTestTableSQL select table 'tableName'
func SelectTestTableWithTableName(db *sqlx.DB, dest *[]TestTable, tableName string, keys ...interface{}) error {
	var tmpKeys []interface{}
	tmpKeys = append(tmpKeys, tableName)
	tmpKeys = append(tmpKeys, keys...)
	sqlStr := fmt.Sprintf(SelectTestTableSQL, tmpKeys...)
	return db.Select(dest, sqlStr)
}

// InsertTestTable use InsertTestTableSQL insert table TestTable
func InsertTestTable(db *sqlx.DB, src TestTable) (sql.Result, error) {
	return InsertTestTableWithTableName(db, TestTableTableName, src)
}

// InsertTestTableWithTableName use InsertTestTableSQL insert table 'tableName'
func InsertTestTableWithTableName(db *sqlx.DB, tableName string, src TestTable) (sql.Result, error) {
	sqlStr := fmt.Sprintf(InsertTestTableSQL, tableName)
	return db.NamedExec(sqlStr, src)
}

// DeleteTestTable use DeleteTestTableSQL delete table TestTable
func DeleteTestTable(db *sqlx.DB, keys ...interface{}) (sql.Result, error) {
	return DeleteTestTableWithTableName(db, TestTableTableName, keys)
}

// DeleteTestTableWithTableName use DeleteTestTableSQL delete table 'tableName'
func DeleteTestTableWithTableName(db *sqlx.DB, tableName string, keys ...interface{}) (sql.Result, error) {
	var tmpKeys []interface{}
	tmpKeys = append(tmpKeys, tableName)
	tmpKeys = append(tmpKeys, keys...)
	sqlStr := fmt.Sprintf(DeleteTestTableSQL, tmpKeys...)
	return db.Exec(sqlStr)
}
