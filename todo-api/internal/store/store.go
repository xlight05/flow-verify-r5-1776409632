package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/asdlc/todo-api/internal/models"
)

var ErrNotFound = errors.New("todo not found")

type Store struct {
	mu       sync.RWMutex
	todos    map[string]models.Todo
	dataFile string
}

func New(dataFile string) (*Store, error) {
	s := &Store{
		todos:    make(map[string]models.Todo),
		dataFile: dataFile,
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	if err := os.MkdirAll(filepath.Dir(s.dataFile), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(s.dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return s.persist()
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	var list []models.Todo
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	for _, t := range list {
		s.todos[t.ID] = t
	}
	return nil
}

func (s *Store) persist() error {
	list := make([]models.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt < list[j].CreatedAt
	})
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.dataFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.dataFile)
}

func (s *Store) List() []models.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt < list[j].CreatedAt
	})
	return list
}

func (s *Store) Create(t models.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.todos[t.ID] = t
	return s.persist()
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return s.persist()
}
