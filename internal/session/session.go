package session

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
)

type Storage struct {
	mux sync.RWMutex
}

const storageFile = "storage/session.bin"

func (s *Storage) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}
	s.mux.Lock()
	defer s.mux.Unlock()

	sessionBytes, err := os.ReadFile(storageFile)
	if err != nil {
		return nil, errors.Wrap(err, "error on reading session")
	}

	return sessionBytes, nil
}

func (s *Storage) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	err := os.WriteFile(storageFile, data, 0600)
	if err != nil {
		return errors.Wrap(err, "error on writing session")
	}

	return nil
}

func Init() error {
	if err := os.MkdirAll(path.Dir(storageFile), 0700); err != nil {
		return errors.Wrap(err, "storage dir create")
	}
	return nil
}
