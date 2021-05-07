package mysql

import "strings"

// TypesMapping type mapping of mysql to go
var TypesMapping = map[string]string{
	// Integer types
	"tinyint":            "int8",
	"tinyint unsigned":   "uint8",
	"smallint":           "int16",
	"smallint unsigned":  "uint16",
	"mediumint":          "int32",
	"mediumint unsigned": "uint32",
	"int":                "int32",
	"int unsigned":       "uint32",
	"integer":            "int32",
	"integer unsigned":   "uint32",
	"bigint":             "int64",
	"bigint unsigned":    "uint64",
	"bit":                "int64", // This is a fuck da type. PHP-DBAL not support it.
	// Floating point types
	"float":   "float32",
	"double":  "float64",
	"decimal": "float64",
	// Character types.
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",
	"blob":       "string",
	"tinyblob":   "string",
	"mediumblob": "string",
	"longblob":   "string",
	// Date types
	"date":      "string",
	"time":      "string",
	"datetime":  "string",
	"timestamp": "string",
	// Binary types
	"binary":    "string",
	"varbinary": "string",
	// bool & enum
	"bool": "bool",
	"enum": "string",
}

// NullAbleType type mapping of go type to sql.null.
var NullAbleType = map[string]string{
	"int":    "sql.NullInt64",
	"float":  "sql.NullFloat64",
	"string": "sql.NullString",
	"time":   "sql.NullTime",
	"bool":   "sql.NullBool",
}

// GetNullAbleType get sql.null type from go type.
func GetNullAbleType(goType string) (nullType string) {
	for k, v := range NullAbleType {
		if strings.Contains(goType, k) {
			nullType = v
			break
		}
	}
	return nullType
}
