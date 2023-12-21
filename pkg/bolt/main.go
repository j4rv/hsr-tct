package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

type database struct {
	db *bolt.DB
}

func New() *database {
	return &database{}
}

func (database *database) Init(databasePath string) (func() error, error) {
	db, err := bolt.Open(databasePath, 0600, nil)
	if err != nil {
		db.Close()
		return nil, err
	}
	database.db = db
	err = database.CreateBuckets()
	return db.Close, err
}

func (database *database) CreateBuckets() error {
	// Start a writable transaction.
	tx, err := database.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists(charactersBucket)
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists(lightconesBucket)
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists(enemiesBucket)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// getEntity retrieves the JSON-encoded data of an entity from the specified BoltDB bucket based on the provided key.
// The data is unmarshaled into the provided type T. If the entity is not found, an error is returned.
func getEntity[T any](db *bolt.DB, bucketName []byte, key string) (T, error) {
	var result T

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %s", string(bucketName))
		}

		jsonData := bucket.Get([]byte(key))
		if jsonData == nil {
			return fmt.Errorf("entity with key '%s' not found", key)
		}

		if err := json.Unmarshal(jsonData, &result); err != nil {
			return fmt.Errorf("error unmarshaling JSON data for key '%s': %v", key, err)
		}

		return nil
	})

	return result, err
}

// addEntity adds an entity to the specified BoltDB bucket with the provided key.
// The entity is first serialized to JSON, and the JSON data is stored in the database.
// If the key already exists, the existing value will be replaced with the new value.
func addEntity[T any](db *bolt.DB, bucket []byte, key string, entity T) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		jsonData, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), jsonData)
		return err
	})
}

// deleteEntity deletes an entity from the specified BoltDB bucket based on the provided key.
// If the entity is not found, an error is returned. The function opens an update transaction
// to delete the key-value pair associated with the provided key.
func deleteEntity(db *bolt.DB, bucketName []byte, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %s", string(bucketName))
		}

		// Check if the entity exists before attempting to delete
		if existingData := bucket.Get([]byte(key)); existingData == nil {
			return fmt.Errorf("entity with key '%s' not found", key)
		}

		return bucket.Delete([]byte(key))
	})
}

func NewUUID() []byte {
	return uuid.NewV4().Bytes()
}
