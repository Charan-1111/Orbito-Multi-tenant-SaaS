package config

import (
	"os"
	"sync"
	"task-manager/models"
	"task-manager/internal/store/database"

	"github.com/bytedance/sonic"
)

type Configuration struct {
	Environment string        `json:"environment"`
	Server      models.Server `json:"server"`
	Database *database.Database `json:"database"`
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
