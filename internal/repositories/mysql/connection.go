package mysql

import (
	"database/sql"
	"log"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New() *Store {
	cfg := mysqldriver.Config{
		User:   "root",
		Passwd: "secret",
		Net:    "tcp",
		Addr:   "127.0.0.1:33070",
		DBName: "test",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	log.Println(err)
	log.Println(db.Ping())

	return &Store{db}
}

// docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=test -p 33070:3306 -d mysql:8.0
