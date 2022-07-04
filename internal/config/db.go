package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func DBNew() *mysql.Config {

	conf := mysql.NewConfig()

	conf.User = os.Getenv("DB_USER")
	conf.Passwd = os.Getenv("MYSQL_ROOT_PASSWORD")
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	conf.DBName = os.Getenv("MYSQL_DATABASE")

	log.Printf("%+v", conf)

	return conf
}
