package database

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	// these are the sensitive information
	// get them from somewhere else for now we can keep them here
	Username     string `json:"username"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName string `json:"databaseName"`
}

type DataBaseStore struct {
	Db   *pgxpool.Pool
	once sync.Once
}

func (db *DataBaseStore) InitializeDatabaseStore(ctx context.Context, database *Database) error {
	if database == nil {
		return errors.New("database configuration is required")
	}

	var err error

	db.once.Do(func() {
		dsn := "postgres://" + database.Username + ":" + database.Password + "@" + database.Host + ":" + database.Port + "/" + database.DatabaseName
		db.Db, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return
		}
	})

	return fmt.Errorf("Error while initializing the database connection pool : %w", err)
}

func (db *DataBaseStore) Close() {
	if db.Db != nil {
		db.Db.Close()
	}
}
