package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

func dbConnect(retriesLeft int) error {
	con := os.Getenv("POSTGRES")
	logger.Infof("opening connection to postgres...")

	var err error
	db, err = sql.Open("postgres", con)
	if err != nil {
		logger.Errorf("cannot open connection to postgres: %v", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		if retriesLeft > 0 {
			logger.Errorf("cannot talk to postgres, retrying in 10 seconds (%d attempts left): %v", retriesLeft-1, err)
			time.Sleep(10 * time.Second)
			return dbConnect(retriesLeft - 1)
		} else {
			logger.Errorf("cannot talk to postgres, last attempt failed: %v", err)
			return err
		}
	}

	statement := `
    CREATE TABLE IF NOT EXISTS migrations (
      filename TEXT NOT NULL UNIQUE
    );
  `
	_, err = db.Exec(statement)
	if err != nil {
		logger.Errorf("cannot create migrations table: %v", err)
		return err
	}

	maxIdleConnections, err := strconv.Atoi(os.Getenv("MAX_IDLE_PG_CONNECTIONS"))
	if err != nil {
		logger.Warningf("cannot parse COMMENTO_MAX_IDLE_PG_CONNECTIONS: %v", err)
		maxIdleConnections = 50
	}

	db.SetMaxIdleConns(maxIdleConnections)

	return nil
}
