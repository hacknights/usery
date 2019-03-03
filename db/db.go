package db

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
)

type DB struct {
	instance *buntdb.DB
}

func NewDB(d *buntdb.DB) *DB {
	return &DB{
		instance: d,
	}
}

func (db *DB) Insert(e entity) (identifier, error) {
	assignIdentity(e)
	if err := db.instance.Update(func(tx *buntdb.Tx) error {
		if _, err := tx.Get(e.Identify()); err == nil {
			return errors.New("")
		}

		s, err := e.serialize()
		if err != nil {
			return err
		}

		if _, _, err := tx.Set(e.Identify(), s, nil); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return e, nil
}

func (db *DB) Find(id string, e entity) error {
	if err := db.instance.View(func(tx *buntdb.Tx) error {
		s, err := tx.Get(id)
		if err != nil {
			return err
		}

		if err := e.deserialize(s); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type entity interface {
	identifier
	identityAssigner
	serialize() (string, error)
	deserialize(string) error
}

type identityAssigner interface {
	identityAssign(id string)
}

type identifier interface {
	Identify() string
}

func assignIdentity(i identityAssigner) {
	i.identityAssign(uuid.New().String())
}
