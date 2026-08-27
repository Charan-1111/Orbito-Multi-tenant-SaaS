package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"task-manager-auth/internal/logging"
	"task-manager-auth/internal/models"

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
	Db      *pgxpool.Pool
	Queries models.Queries
	Log     *logging.Log
	once    sync.Once
}

func (db *DataBaseStore) InitializeDatabaseStore(ctx context.Context, database *Database, queries models.Queries) error {
	if database == nil {
		return errors.New("database configuration is required")
	}

	db.Queries = queries

	var err error

	db.once.Do(func() {
		dsn := "postgres://" + database.Username + ":" + database.Password + "@" + database.Host + ":" + database.Port + "/" + database.DatabaseName
		db.Db, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return
		}
	})

	if err != nil {
		return fmt.Errorf("Error initializing database connection pool : %w", err)
	}

	return nil
}

func (db *DataBaseStore) PingDb(ctx context.Context) error {
	return db.Db.Ping(ctx)
}

func (db *DataBaseStore) Close() {
	if db.Db != nil {
		db.Db.Close()
	}
}
