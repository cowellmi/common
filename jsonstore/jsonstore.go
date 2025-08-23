package jsonstore

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	dir  string
	path string
}

func NewStore(path string) (*Store, error) {
	dir, fName := filepath.Split(path)
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if fName == "" {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrInvalid}
	}

	return &Store{
		dir:  dir,
		path: path,
	}, nil
}

func (s *Store) Read(v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(b) == 0 {
		return nil
	}

	return json.Unmarshal(b, v)
}

func (s *Store) Write(v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(s.dir, "")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()
	didRename := false
	defer func() {
		tmp.Close()
		if !didRename {
			os.Remove(tmpName)
		}
	}()

	_, err = tmp.Write(b)
	if err != nil {
		return err
	}

	err = tmp.Sync()
	if err != nil {
		return err
	}

	err = tmp.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tmpName, s.path)
	if err != nil {
		return err
	}

	didRename = true
	return nil
}
