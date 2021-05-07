package generator

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"os/exec"
	"zgenerator/config"
	"zgenerator/generator/mysql"
)

// Generate generate go code in different sql driver
func Generate() error {
	cfg := config.GetGlobalConfig()
	db := sqlx.MustConnect(cfg.DriverName, cfg.DSN)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Warning: DB connection close failed.", err)
		}
	}(db)
	err := db.Ping()
	if err != nil {
		return err
	}
	var code string
	switch cfg.DriverName {
	case "mysql":
		code, err = mysql.GenerateCode(db)
	case "postgres":
		err = errors.New("TODO: postgresql driver")
	case "sqlite3":
		err = errors.New("TODO: sqlite3 driver")
	case "mssql":
		err = errors.New("TODO: mssql driver")
	default:
		err = errors.New("unknown driver type")
	}
	if err == nil {
		err = WriteCodeToFile(cfg, code)
	}
	if err == nil {
		cmd := exec.Command("gofmt", "-w", cfg.FilePath)
		err = cmd.Run()
	}
	return err
}

// WriteCodeToFile write code to target file.(TRUNC)
func WriteCodeToFile(cfg *config.Config, code string) (err error) {
	file, err := os.OpenFile(cfg.FilePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("File close failed.", err)
		}
	}(file)
	_, err = file.WriteString(code)
	if err != nil {
		return
	}
	return nil
}
