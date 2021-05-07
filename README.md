# zgenerator
Generate Golang code according to the database table structure.

## Quick start

```shell
./zgenerator -dsn="root:123456@tcp(127.0.0.1:3306)/TestDB?charset=utf8mb4" -table-name="TestTable" -output="./MyTest.go" -first-to-upper=true
```

Table create SQL:
```sql
CREATE Table TestTable (ID BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'UserID', Name VARCHAR(30) COMMENT 'User Name', Money BIGINT COMMENT 'Money');
```

You can get Golang file like this:
```go
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

```

## Usage
```shell
Usage of zgenerator:
  -charset string
        Charset of db connection (default "utf8mb4")
  -code-after-struct string
        Customize code after struct
  -code-before-struct string
        Customize code before struct
  -code-in-struct string
        Customize code in struct
  -date-string
        Use string store date type (default true)
  -db string
        Database name
  -db-driver string
        Driver of database (default "mysql")
  -dsn string
        Data Source Name
  -first-to-upper
        Convert first letter to upper
  -generate-func
        Generate some easy function to use (Depend sqlx) (default true)
  -help
        Show usage
  -host string
        Host of database (default "127.0.0.1")
  -json-tag
        Add json tag
  -output string
        Path of output file
  -package-name string
        package name of generate code (default "handler")
  -password string
        Password of user
  -post uint
        Port of database (default 3306)
  -select-key string
        Key of select SQL (Default: PrimaryKey)
  -struct-name string
        Struct name of generate code (Default: same with table-name)
  -table-name string
        Table name in database
  -user string
        User of database (default "root")
  -version
        Version info
```