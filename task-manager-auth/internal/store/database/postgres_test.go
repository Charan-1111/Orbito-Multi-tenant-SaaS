package database

import (
	"context"
	"testing"
)

func TestInitializeDatabaseStoreRequiresConfiguration(t *testing.T) {
	store := &DataBaseStore{}
	if err := store.InitializeDatabaseStore(context.Background(), nil); err == nil {
		t.Fatal("expected an error for nil database configuration")
	}
}

func TestCloseAllowsUninitializedStore(t *testing.T) {
	store := &DataBaseStore{}
	store.Close()
}
