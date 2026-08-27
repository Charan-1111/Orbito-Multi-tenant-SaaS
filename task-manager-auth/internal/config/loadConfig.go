package config

import (
	"os"
	"sync"
	"task-manager-auth/internal/models"
	"task-manager-auth/internal/store/database"

	"github.com/bytedance/sonic"
)

type Configuration struct {
	Environment string             `json:"environment"`
	Server      models.Server      `json:"server"`
	Database    *database.Database `json:"database"`
	Queries     models.Queries     `json:"queries"`
	once        sync.Once
}

func (config *Configuration) LoadConfig() error {
	filePath := os.Getenv("CONFIG_FILE_PATH")
	if filePath == "" {
		filePath = "configs/local/config.json"
	}

	var loadErr error

	config.once.Do(func() {
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			loadErr = err
			return
		}

		err = sonic.Unmarshal(fileBytes, &config)
		if err != nil {
			loadErr = err
			return
		}
	})

	return loadErr
}
