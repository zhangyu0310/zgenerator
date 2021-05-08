package mysql

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"zgenerator/config"
	"zgenerator/generator/util"
)

// Column is struct of MySQL INFORMATION_SCHEMA.COLUMNS table
type Column struct {
	TableSchema     string `db:"TABLE_SCHEMA"`
	TableName       string `db:"TABLE_NAME"`
	ColumnName      string `db:"COLUMN_NAME"`
	OrdinalPosition uint64 `db:"ORDINAL_POSITION"`
	IsNullAble      string `db:"IS_NULLABLE"`
	DataType        string `db:"DATA_TYPE"`
	ColumnType      string `db:"COLUMN_TYPE"`
	ColumnKey       string `db:"COLUMN_KEY"`
	Extra           string `db:"EXTRA"`
	ColumnComment   string `db:"COLUMN_COMMENT"`
}

// GenerateCode generate go code from MySQL table
func GenerateCode(db *sqlx.DB) (code string, err error) {
	cfg := config.GetGlobalConfig()
	// SQL to get INFORMATION_SCHEMA.COLUMNS
	var sqlStr = `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME,
                  ORDINAL_POSITION, IS_NULLABLE, DATA_TYPE,
                  COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT
                  FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE()
                  AND TABLE_NAME='%s' ORDER BY ORDINAL_POSITION asc`
	sqlStr = fmt.Sprintf(sqlStr, cfg.TableName)
	var columns []Column
	err = db.Select(&columns, sqlStr)
	if err != nil {
		return
	}
	code, err = getCode(cfg, columns)
	if err != nil {
		return
	}
	return
}

// getCode real function for generate go code
// TODO: 拆分函数，添加注释，整理逻辑。这个函数太乱了。。。
func getCode(cfg *config.Config, columns []Column) (code string, err error) {
	code = fmt.Sprintf("// Package %s is generate by zgenerator (Author:poppinzhang).\n", cfg.PackageName)
	code += fmt.Sprintf("package %s\n\n", cfg.PackageName)
	if !cfg.DateString {
		TypesMapping["date"] = "time.Time"
		TypesMapping["time"] = "time.Time"
		TypesMapping["datetime"] = "time.Time"
		TypesMapping["timestamp"] = "time.Time"
		code += "import \"time\"\n\n"
	}
	if cfg.GenerateFunc {
		code += "import \"database/sql\"\n"
		code += "import \"fmt\"\n"
		code += "import \"github.com/jmoiron/sqlx\"\n\n"
	}
	if cfg.CustomizeCodeBefore != "" {
		code += fmt.Sprintf("%s\n\n", cfg.CustomizeCodeBefore)
	}
	structName := cfg.StructName
	if cfg.FirstToUpper {
		structName = util.FirstToUpper(structName)
	}
	code += fmt.Sprintf("// %s from Database %s, Table %s\n",
		structName, columns[0].TableSchema, cfg.TableName)
	code += fmt.Sprintf("type %s struct {\n", structName)

	// Loop processing columns
	//TODO: 封装成函数    生成struct
	insertSQL1 := fmt.Sprintf("var Insert%sSQL = \"INSERT INTO `%%v` (", structName)
	insertSQL2 := ") VALUES ("
	insertSQL3 := ") ON DUPLICATE KEY UPDATE "
	var pkVec []string
	for _, column := range columns {
		goType, ok := TypesMapping[column.DataType]
		if !ok {
			err = errors.New("not support mysql type " + column.DataType)
			return
		}
		if column.IsNullAble == "YES" {
			goType = GetNullAbleType(goType)
		} else {
			if strings.Contains(strings.ToLower(column.ColumnType), "unsigned") {
				tmpType := column.DataType + " unsigned"
				goType, ok = TypesMapping[tmpType]
				if !ok {
					err = errors.New("not support mysql type " + tmpType)
					return
				}
			}
		}
		columnName := column.ColumnName
		if cfg.FirstToUpper {
			columnName = util.FirstToUpper(columnName)
		}
		if cfg.JsonTag {
			code += fmt.Sprintf("%s %s `db:\"%s\" json:\"%s\"`",
				columnName, goType, column.ColumnName, column.ColumnName)
		} else {
			code += fmt.Sprintf("%s %s `db:\"%s\"`",
				columnName, goType, column.ColumnName)
		}
		if column.ColumnComment != "" {
			code += fmt.Sprintf(" // %s", column.ColumnComment)
		}
		code += "\n"

		if strings.ToUpper(column.Extra) != "AUTO_INCREMENT" {
			insertSQL1 += fmt.Sprintf("`%s`, ", column.ColumnName)
			insertSQL2 += fmt.Sprintf(":%s, ", columnName)
		}
		if column.ColumnKey == "PRI" || column.ColumnKey == "UNI" {
			if column.ColumnKey == "PRI" {
				pkVec = append(pkVec, column.ColumnName)
			}
		} else {
			if strings.ToUpper(column.Extra) != "AUTO_INCREMENT" {
				insertSQL3 += fmt.Sprintf("`%s`=:%s, ", column.ColumnName, columnName)
			}
		}
	}
	insertSQL1 = strings.TrimSuffix(insertSQL1, ", ")
	insertSQL2 = strings.TrimSuffix(insertSQL2, ", ")
	insertSQL3 = strings.TrimSuffix(insertSQL3, ", ")
	insertSQL := insertSQL1 + insertSQL2 + insertSQL3

	if cfg.CustomizeCodeIn != "" {
		code += fmt.Sprintf("%s\n", cfg.CustomizeCodeIn)
	}
	code += "}\n\n"
	if cfg.CustomizeCodeAfter != "" {
		code += fmt.Sprintf("%s\n\n", cfg.CustomizeCodeAfter)
	}

	// TODO: 封装成函数  生成操作SQL
	selectSQL := fmt.Sprintf("var Select%sSQL = \"SELECT * FROM `%%v` WHERE ", structName)
	deleteSQL := fmt.Sprintf("var Delete%sSQL = \"DELETE FROM `%%v` WHERE ", structName)
	if cfg.SelectKey != "" {
		selectSQL += fmt.Sprintf("`%s`='%%v'\"\n\n", cfg.SelectKey)
		deleteSQL += fmt.Sprintf("`%s`='%%v'\"\n\n", cfg.SelectKey)
	} else {
		for _, pk := range pkVec {
			selectSQL += fmt.Sprintf("`%s`='%%v' AND ", pk)
			deleteSQL += fmt.Sprintf("`%s`='%%v' AND ", pk)
		}
		selectSQL = strings.TrimSuffix(selectSQL, " AND ")
		deleteSQL = strings.TrimSuffix(deleteSQL, " AND ")
	}

	code += fmt.Sprintf("// Select%sSQL select SQL for %s\n", structName, structName)
	code += fmt.Sprintf("%s\"\n\n", selectSQL)
	code += fmt.Sprintf("// Insert%sSQL insert SQL for %s\n", structName, structName)
	code += fmt.Sprintf("%s\"\n\n", insertSQL)
	code += fmt.Sprintf("// Delete%sSQL delete SQL for %s\n", structName, structName)
	code += fmt.Sprintf("%s\"\n\n", deleteSQL)

	code += fmt.Sprintf("// %sTableName table name from user input.\n", structName)
	code += fmt.Sprintf("var %sTableName = \"%s\"\n\n", structName, cfg.TableName)

	// TODO: 封装成函数   生成操作函数
	if cfg.GenerateFunc {
		code += fmt.Sprintf("// Select%s use Select%sSQL select table %s\n", structName, structName, cfg.TableName)
		code += fmt.Sprintf("func Select%s(db *sqlx.DB, dest *[]%s, keys ...interface{}) error {\n", structName, structName)
		code += fmt.Sprintf("return Select%sWithTableName(db, dest, %sTableName, keys)}\n\n", structName, structName)

		code += fmt.Sprintf("// Select%sWithTableName use Select%sSQL select table 'tableName'\n", structName, structName)
		code += fmt.Sprintf("func Select%sWithTableName(db *sqlx.DB, dest *[]%s, tableName string, keys ...interface{}) error {\n", structName, structName)
		code += fmt.Sprintf("var tmpKeys []interface{}\ntmpKeys = append(tmpKeys, tableName)\ntmpKeys = append(tmpKeys, keys...)\n")
		code += fmt.Sprintf("sqlStr := fmt.Sprintf(Select%sSQL, tmpKeys...)\n", structName)
		code += fmt.Sprintf("return db.Select(dest, sqlStr)}\n\n")

		code += fmt.Sprintf("// Insert%s use Insert%sSQL insert table %s\n", structName, structName, cfg.TableName)
		code += fmt.Sprintf("func Insert%s(db *sqlx.DB, src %s) (sql.Result, error) {\n", structName, structName)
		code += fmt.Sprintf("return Insert%sWithTableName(db, %sTableName, src)}\n\n", structName, structName)

		code += fmt.Sprintf("// Insert%sWithTableName use Insert%sSQL insert table 'tableName'\n", structName, structName)
		code += fmt.Sprintf("func Insert%sWithTableName(db *sqlx.DB, tableName string, src %s) (sql.Result, error) {\n", structName, structName)
		code += fmt.Sprintf("sqlStr := fmt.Sprintf(Insert%sSQL, tableName)\nreturn db.NamedExec(sqlStr, src)\n}\n\n", structName)

		code += fmt.Sprintf("// Delete%s use Delete%sSQL delete table %s\n", structName, structName, cfg.TableName)
		code += fmt.Sprintf("func Delete%s(db *sqlx.DB, keys ...interface{}) (sql.Result, error) {\n", structName)
		code += fmt.Sprintf("\treturn Delete%sWithTableName(db, %sTableName, keys)\n}\n\n", structName, structName)

		code += fmt.Sprintf("// Delete%sWithTableName use Delete%sSQL delete table 'tableName'\n", structName, structName)
		code += fmt.Sprintf("func Delete%sWithTableName(db *sqlx.DB, tableName string, keys ...interface{}) (sql.Result, error) {\n", structName)
		code += fmt.Sprintf("var tmpKeys []interface{}\ntmpKeys = append(tmpKeys, tableName)\ntmpKeys = append(tmpKeys, keys...)\n")
		code += fmt.Sprintf("sqlStr := fmt.Sprintf(Delete%sSQL, tmpKeys...)\n\treturn db.Exec(sqlStr)}\n", structName)
	}
	//TODO: 位置整理  如果没有使用到time.Time 类型，删除import time
	if !cfg.DateString && !strings.Contains(code, "time.Time") {
		code = strings.Replace(code, "import \"time\"\n\n", "", 1)
	}
	return
}
