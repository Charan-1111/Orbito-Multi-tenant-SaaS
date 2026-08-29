package database

import (
	"context"
	"task-manager-auth/internal/models"
	"testing"
)

func TestInitializeDatabaseStoreRequiresConfiguration(t *testing.T) {
	store := &DataBaseStore{}
	if err := store.InitializeDatabaseStore(context.Background(), models.Queries{}); err == nil {
		t.Fatal("expected an error for nil database configuration")
	}
}

func TestCloseAllowsUninitializedStore(t *testing.T) {
	store := &DataBaseStore{}
	store.Close()
}
