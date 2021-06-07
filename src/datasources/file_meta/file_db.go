package file_meta

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	mysqlUsername = "mysql_username"
	mysqlPassword = "mysql_password"
	mysqlHost     = "mysql_host"
	mysqlSchema   = "mysql_schema"
)

var (
	Client *sql.DB

	username = os.Getenv(mysqlUsername)
	password = os.Getenv(mysqlPassword)
	host     = os.Getenv(mysqlHost)
	schema   = os.Getenv(mysqlSchema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)
	log.Println("datasource= ", dataSourceName)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}
