package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"social-network/internal/errapp"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New(cfg mysqldriver.Config) (*Store, error) {
	dsn := cfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %v; DSN: %s", err, dsn)
	}

	err = tryConnect(db)
	if err != nil {
		return nil, fmt.Errorf("tryConnect error: %v", err)
	}

	err = initTables(db)
	if err != nil {
		return nil, fmt.Errorf("initTables error: %v", err)
	}

	store := Store{
		db: db,
	}

	return &store, nil
}

func tryConnect(db *sql.DB) error {
	numAttempt := 1

	for db.Ping() != nil {

		if numAttempt > 10 {
			log.Printf("fail db ping after several attempts")
			return errapp.DBPing
		}

		sleepTime := time.Duration(numAttempt) * time.Second
		log.Printf("fail db ping; attempt %v; need sleep %s seconds", numAttempt, sleepTime)
		time.Sleep(sleepTime)

		numAttempt++
	}

	log.Printf("success db ping; attempt %v", numAttempt)
	return nil
}

func initTables(db *sql.DB) (err error) {
	_, err = db.Exec(`
		create table if not exists user_data
		(
			id      varchar(36)   not null
				primary key,
			name    varchar(100)  null,
			surname varchar(100)  null,
			age     smallint      null,
			gender  char          null,
			hobbies varchar(1000) null,
			city    varchar(100)  null
		);
	`)

	if err != nil {
		return fmt.Errorf("create table user_data error: %v", err)
	}

	_, err = db.Exec(`
		create table if not exists user_access
		(
			login    varchar(20)  not null
				primary key,
			password varchar(100) not null,
			user_id  varchar(36)  null,
			constraint user_access_user_data_id_fk
				foreign key (user_id) references user_data (id)
					on update cascade on delete cascade
		);
	`)

	if err != nil {
		return fmt.Errorf("create table user_access error: %v", err)
	}

	_, err = db.Exec(`
		create table if not exists friends
		(
			user_id   varchar(36) not null,
			friend_id varchar(36) not null,
			primary key (user_id, friend_id),
			constraint friend_id_fk
				foreign key (friend_id) references user_data (id)
					on update cascade on delete cascade,
			constraint user_id_fk
				foreign key (user_id) references user_data (id)
					on update cascade on delete cascade
		);
	`)

	if err != nil {
		return fmt.Errorf("create table friends error: %v", err)
	}

	return nil
}
