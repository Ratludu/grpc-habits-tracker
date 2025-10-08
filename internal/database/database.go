package database

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
)

type Database struct {
	db *badger.DB
}

func (d *Database) Get(key []byte) ([]byte, error) {
	var valCopy []byte

	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err // Returns badger.ErrKeyNotFound if key is not found.
		}

		valCopy, err = item.ValueCopy(nil)
		return err // Returns any error encountered during value copy.
	})

	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("key not found: %s", string(key))
	}

	if err != nil {
		return nil, fmt.Errorf("error viewing database: %w", err)
	}

	return valCopy, nil
}

func (d *Database) GetAll() (map[string]string, error) {

	resultMap := make(map[string]string)

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			resultMap[string(k)] = string(val)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

func Open(dbPath string) (*Database, error) {
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB at %s: %w", dbPath, err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Set(key, value []byte) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)

		return txn.SetEntry(entry)
	})

	if err != nil {
		return fmt.Errorf("failed to set cache value: %w", err)
	}
	return nil
}

func (d *Database) Delete(key []byte) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
