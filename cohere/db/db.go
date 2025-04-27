package db

import (
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v4"
)

type Database struct {
	db     *badger.DB
	dbPath string
}

// NewDatabase returns a Database struct containing a BadgerDB instance
func NewDatabase(path string) (*Database, error) {
	badgerDb, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open Badger database at %s: %v", path, err)
	}

	db := &Database{
		db:     badgerDb,
		dbPath: path,
	}

	return db, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Cleanup closes the database connection and removes the database directory
func (d *Database) Cleanup() error {
	if err := d.Close(); err != nil {
		return fmt.Errorf("failed to close database during cleanup: %v", err)
	}

	if err := os.RemoveAll(d.dbPath); err != nil {
		return fmt.Errorf("failed to remove database directory during cleanup: %v", err)
	}

	return nil
}

// GetKey retrieves the value associated with the provided key from the BadgerDB instance.
func (d *Database) GetKey(key string) ([]byte, error) {
	var valCopy []byte
	err := d.db.View(func(txn *badger.Txn) error {
		val, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("key '%s' not found", key)
			}
			return fmt.Errorf("failed to get key '%s': %v", key, err)
		}
		valCopy, err = val.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("failed to copy value for key '%s': %v", key, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed while getting key '%s': %v", key, err)
	}

	return valCopy, nil
}

// SetKey sets the value associated with the provided key in the BadgerDB instance
func (d *Database) SetKey(key string, value string) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
		if err := txn.SetEntry(e); err != nil {
			return fmt.Errorf("failed to set key '%s' with value '%s': %v", key, value, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed while setting key '%s': %v", key, err)
	}

	return nil
}

// DeleteKey removes the key-value pair associated with the provided key from the BadgerDB instance
func (d *Database) DeleteKey(key string) error {
	_, err := d.GetKey(key)
	if err == badger.ErrKeyNotFound {
		return fmt.Errorf("key '%s' not found, cannot delete", key)
	} else if err != nil {
		return fmt.Errorf("failed to check existence of key '%s': %v", key, err)
	}

	err = d.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(key)); err != nil {
			return fmt.Errorf("failed to delete key '%s': %v", key, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed while deleting key '%s': %v", key, err)
	}

	return nil
}

func (d *Database) GetKeys() ([][]byte, error) {
	var keys [][]byte

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		iter := txn.NewIterator(opts)
		defer iter.Close()

		for iter.Rewind(); iter.Valid(); iter.Next() {
			key := iter.Item().Key()
			keys = append(keys, key)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed while getting keys: %v", err)
	}

	return keys, nil
}
