package sfjdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type DB[T any] struct {
	rw       sync.RWMutex
	data     T
	filepath string
}

// Save saves a copy of the data as plain json in the file.
func (db *DB[T]) Save(data T) error {
	db.rw.Lock()
	defer db.rw.Unlock()
	db.data = *objcopy[T](data)
	content, err := json.Marshal(db.data)
	if err != nil {
		return err
	}
	return WriteFile(db.filepath, content, 0644)
}

// Read returns a copy of the data.
func (db *DB[T]) Read() T {
	db.rw.RLock()
	defer db.rw.RUnlock()
	return *objcopy(db.data)
}

// Filepath returns location of the json file.
func (db *DB[T]) Filepath() string {
	return db.filepath
}

func objcopy[T any](obj T) *T {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	newobj := new(T)
	if err := json.Unmarshal(data, newobj); err != nil {
		panic(err)
	}
	return newobj
}

// WriteFile writes data to filename+some suffix, then renames it into filename.
// The perm argument is ignored on Windows. If the target filename already
// exists but is not a regular file, WriteFile returns an error.
// Copied from:
// https://github.com/tailscale/tailscale/blob/main/atomicfile/atomicfile.go
func WriteFile(filename string, data []byte, perm os.FileMode) (err error) {
	fi, err := os.Stat(filename)
	if err == nil && !fi.Mode().IsRegular() {
		return fmt.Errorf("%s already exists and is not a regular file", filename)
	}
	f, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(tmpName)
		}
	}()
	if _, err := f.Write(data); err != nil {
		return err
	}
	if runtime.GOOS != "windows" {
		if err := f.Chmod(perm); err != nil {
			return err
		}
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, filename)
}

// Open opens a json file as a database.
func Open[T any](filepath string) (db *DB[T], err error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	db = &DB[T]{filepath: filepath}
	if err := json.Unmarshal(content, &db.data); err != nil {
		return nil, err
	}
	return db, nil
}
