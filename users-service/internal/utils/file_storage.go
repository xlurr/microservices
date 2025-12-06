package utils

import (
	"encoding/json"
	"os"
	"sync"
)

// FileStorage - утилита для работы с JSON-хранилищем на диске
type FileStorage struct {
	mu       sync.RWMutex
	filePath string
}

// NewFileStorage создаёт новое хранилище
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

// EnsureFile гарантирует, что файл существует
func (fs *FileStorage) EnsureFile() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		// Создаём файл с пустым массивом
		return os.WriteFile(fs.filePath, []byte("[]"), 0644)
	}
	return nil
}

// LoadJSON загружает данные из JSON файла
func (fs *FileStorage) LoadJSON(v interface{}) error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // файл не существует - это норма
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, v)
}

// SaveJSON сохраняет данные в JSON файл
func (fs *FileStorage) SaveJSON(v interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filePath, data, 0644)
}
