package config

import (
	"github.com/go-sql-driver/mysql"
)

func NewDB() *mysql.Config {
	conf := mysql.NewConfig()

	conf.User = "root"
	conf.Passwd = "pass"
	conf.Net = "tcp"
	conf.Addr = "db:3306"
	conf.DBName = "db"

	return conf
}
