package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New(cfg *mysqldriver.Config) (*Store, error) {
	dsn := cfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %v; DSN: %s", err, dsn)
	}

	err = tryConnect(db)
	if err != nil {
		return nil, fmt.Errorf("tryConnect error: %v; DSN: %s", err, dsn)
	}

	store := Store{
		db: db,
	}

	return &store, nil
}

func tryConnect(db *sql.DB) error {
	numAttempt := 1

	for {
		if err := db.Ping(); err != nil {
			if numAttempt > 10 {
				return errors.New("fail db ping after several attempts")
			}

			sleepTime := time.Duration(numAttempt) * time.Second
			log.Printf("fail db ping: %v; attempt %v; need sleep %s seconds", err, numAttempt, sleepTime)
			time.Sleep(sleepTime)

			numAttempt++
		} else {
			break
		}
	}

	log.Printf("success db ping; attempt %v", numAttempt)
	return nil
}

func transactionRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		if !errors.Is(err, sql.ErrTxDone) {
			log.Printf("tx.Rollback error: %v", err)
		}
	}
}

func rowsClose(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Printf("rows.Close error: %v", err)
	}
}
