package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	DBConnect *sql.DB
}

func NewConnection(driverName, dataSourceName string) (Connection, error) {
	var err error
	var con Connection
	con.DBConnect, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return Connection{}, err
	}
	con.DBConnect.SetConnMaxLifetime(60 * time.Second)
	con.DBConnect.SetMaxIdleConns(6)
	con.DBConnect.SetMaxOpenConns(6)
	return con, err
}
