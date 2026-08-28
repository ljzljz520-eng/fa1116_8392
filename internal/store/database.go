package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"sync"
	"time"
)

var buckets = [][]byte{[]byte("records"), []byte("audits"), []byte("workflows"), []byte("attachments"), []byte("meta")}

type Database struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Database, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepathDir(path), 0o755); err != nil {
		return nil, fmt.Errorf("prepare database directory: %w", err)
	}
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	d := &Database{db: db}
	if err := d.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}

func filepathDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i == 0 {
				return "/"
			}
			return path[:i]
		}
	}
	return "."
}

func (d *Database) initialize() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.db == nil {
		return nil
	}
	err := d.db.Close()
	d.db = nil
	return err
}

func (d *Database) Put(bucket, key string, value any) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return errors.New("database is closed")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return b.Put([]byte(key), data)
	})
}

func (d *Database) Get(bucket, key string, target any) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return false, errors.New("database is closed")
	}
	found := false
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		data := b.Get([]byte(key))
		if data == nil {
			return nil
		}
		found = true
		return json.Unmarshal(data, target)
	})
	return found, err
}

func (d *Database) Delete(bucket, key string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return b.Delete([]byte(key))
	})
}

func (d *Database) List(bucket string) ([][]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return nil, errors.New("database is closed")
	}
	var values [][]byte
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return b.ForEach(func(_, v []byte) error {
			if v != nil {
				values = append(values, append([]byte(nil), v...))
			}
			return nil
		})
	})
	return values, err
}

func (d *Database) Count(bucket string) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return 0, errors.New("database is closed")
	}
	count := 0
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}
