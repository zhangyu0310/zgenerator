package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"zgenerator/config"
	"zgenerator/generator"
)

var (
	dsn                 = flag.String("dsn", "", "Data Source Name")
	user                = flag.String("user", "root", "User of database")
	passwd              = flag.String("password", "", "Password of user")
	host                = flag.String("host", "127.0.0.1", "Host of database")
	port                = flag.Uint("post", 3306, "Port of database")
	dbName              = flag.String("db", "", "Database name")
	charset             = flag.String("charset", "utf8mb4", "Charset of db connection")
	driverName          = flag.String("db-driver", "mysql", "Driver of database")
	tableName           = flag.String("table-name", "", "Table name in database")
	structName          = flag.String("struct-name", "", "Struct name of generate code (Default: same with table-name)")
	filePath            = flag.String("output", "", "Path of output file")
	packageName         = flag.String("package-name", "handler", "package name of generate code")
	jsonTag             = flag.Bool("json-tag", false, "Add json tag")
	dateString          = flag.Bool("date-string", true, "Use string store date type")
	selectKeys          = flag.String("select-keys", "", "Key of select SQL (Use delimiter ',' | Default: PrimaryKey)")
	customizeCodeBefore = flag.String("code-before-struct", "", "Customize code before struct")
	customizeCodeIn     = flag.String("code-in-struct", "", "Customize code in struct")
	customizeCodeAfter  = flag.String("code-after-struct", "", "Customize code after struct")
	firstToUpper        = flag.Bool("first-to-upper", false, "Convert first letter to upper")
	generateFunc        = flag.Bool("generate-func", true, "Generate some easy function to use (Depend sqlx)")
)

// checkCmd check if input commands are valid.
func checkCmd() error {
	if *dsn == "" {
		if *dbName == "" {
			return errors.New("DB name is required")
		}
		*dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
			*user, *passwd, *host, *port, *dbName, *charset)
	}
	fmt.Println("DSN is", *dsn)
	db, err := sql.Open(*driverName, *dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	if *tableName == "" {
		return errors.New("table name is required")
	}
	if *structName == "" {
		*structName = *tableName
	}
	if *filePath == "" {
		*filePath = "./" + *tableName + ".go"
	}
	return nil
}

// cmdConfigSetToGlobal store command config to global config.
func cmdConfigSetToGlobal(cfg *config.Config) {
	cfg.DSN = *dsn
	cfg.DriverName = *driverName
	cfg.TableName = *tableName
	cfg.StructName = *structName
	cfg.FilePath = *filePath
	cfg.PackageName = *packageName
	cfg.JsonTag = *jsonTag
	cfg.DateString = *dateString
	cfg.SelectKeys = *selectKeys
	cfg.CustomizeCodeBefore = *customizeCodeBefore
	cfg.CustomizeCodeIn = *customizeCodeIn
	cfg.CustomizeCodeAfter = *customizeCodeAfter
	cfg.FirstToUpper = *firstToUpper
	cfg.GenerateFunc = *generateFunc
}

func main() {
	help := flag.Bool("help", false, "Show usage")
	ver := flag.Bool("version", false, "Version info")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *ver {
		// FIXME: Print version info.(Version/GitCommit/CompileTime...)
		fmt.Println("Version: v0.0.1")
		os.Exit(0)
	}
	err := checkCmd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.InitializeConfig(cmdConfigSetToGlobal)
	err = generator.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
