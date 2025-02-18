package users

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

const usersFileName = "users.json"

type Storage struct {
	dataDirectory string
	mu            sync.RWMutex
}

func NewStorage(dataDirectory string) (*Storage, error) {
	storage := &Storage{dataDirectory: dataDirectory}
	if err := storage.ensureUsersFileExists(); err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *Storage) Save(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, err := s.readAllUsersFromFile()
	if err != nil {
		return err
	}

	updated := false
	for i, u := range users {
		if u.Username == user.Username {
			users[i] = *user
			updated = true
			break
		}
	}

	if !updated {
		users = append(users, *user)
	}

	return s.overwriteUsersFile(users)
}

func (s *Storage) Get(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users, err := s.readAllUsersFromFile()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, nil
}

func (s *Storage) GetAll() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.readAllUsersFromFile()
}

func (s *Storage) readAllUsersFromFile() ([]User, error) {
	filePath := fmt.Sprintf("%s/%s", s.dataDirectory, usersFileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open users file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

func (s *Storage) overwriteUsersFile(users []User) error {
	filePath := fmt.Sprintf("%s/%s", s.dataDirectory, usersFileName)

	bytes, err := json.MarshalIndent(users, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

func (s *Storage) ensureUsersFileExists() error {
	newFileContent := "[]"

	filePath := fmt.Sprintf("%s/%s", s.dataDirectory, usersFileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(newFileContent), 0644); err != nil {
			return fmt.Errorf("failed to create users file: %w", err)
		}
	}

	return nil
}
