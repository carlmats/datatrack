package database

import (
	"datatrack/database/ephemeral"
	"datatrack/model"
	"errors"

	"github.com/boltdb/bolt"
)

// SetUser sets the user info, replacing if already in place.
func SetUser(u model.User) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return err
		}

		err = db.Put([]byte("name"), ephemeral.Encrypt([]byte(u.Name)))
		if err != nil {
			return err
		}

		return db.Put([]byte("picture"), ephemeral.Encrypt([]byte(u.Picture)))
	})
}

// GetUser returns the user data.
func GetUser() (u *model.User, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte("user"))
		if db == nil {
			return errors.New("no user bucket")
		}

		name := ephemeral.Decrypt(db.Get([]byte("name")))
		if name == nil {
			return errors.New("no name set")
		}
		picture := ephemeral.Decrypt(db.Get([]byte("picture")))
		if name == nil {
			return errors.New("no picture set")
		}

		u = new(model.User)
		u.Name = string(name)
		u.Picture = string(picture)

		return nil
	})
	return
}
