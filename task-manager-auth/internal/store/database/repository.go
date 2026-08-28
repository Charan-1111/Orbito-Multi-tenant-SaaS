package database

import (
	"context"
	"fmt"
)

type Repository interface {
	Close()
	CreateTables(ctx context.Context) error
	RegisterUser(ctx context.Context, username, password string) error
	CheckUserExistance(ctx context.Context, username, passsword string) (int, error)
}

func (db *DataBaseStore) CreateTables(ctx context.Context) error {
	for tableName, query := range db.Queries.Create {
		_, err := db.Db.Exec(ctx, query)
		if err != nil {
			return err
		} else {
			db.Log.Log.Info().Msg("Table " + tableName + " create/already exists")
		}
	}

	return nil
}

func (db *DataBaseStore) RegisterUser(ctx context.Context, username, password string) error {
	_, err := db.Db.Exec(ctx, db.Queries.Insert.RegisterUser, username, password)
	if err != nil {
		return fmt.Errorf("Error while registering the user : w", err)
	}

	return nil
}

func (db *DataBaseStore) CheckUserExistance(ctx context.Context, username, password string) (int, error) {
	var userCnt int

	err := db.Db.QueryRow(ctx, db.Queries.Fetch.CheckUserExistance, username, password).Scan(&userCnt)
	if err != nil {
		return 0, err
	}

	return userCnt, nil
}
