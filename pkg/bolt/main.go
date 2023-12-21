package bolt

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
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

// getEntity retrieves an entity with the specified key from the given BoltDB bucket.
// The entity is unmarshaled from JSON and returned. If the key or bucket does not exist,
// an error is returned.
func getAllEntities[T any](db *bolt.DB, bucketName []byte) ([]T, error) {
	var results []T

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %s", string(bucketName))
		}

		return bucket.ForEach(func(key, value []byte) error {
			var entity T

			if err := json.Unmarshal(value, &entity); err != nil {
				return fmt.Errorf("error unmarshaling JSON data for key '%s': %v", key, err)
			}

			results = append(results, entity)
			return nil
		})
	})

	return results, err
}

// getEntitiesPage retrieves a page of entities from a BoltDB bucket with pagination support.
// The offset parameter specifies the number of records to skip, and the limit parameter specifies
// the maximum number of records to retrieve. If limit is 0, it retrieves all available records starting from the offset.
func getEntitiesPage[T any](db *bolt.DB, bucketName []byte, offset, limit int) ([]T, error) {
	var results []T

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %s", string(bucketName))
		}

		cursor := bucket.Cursor()

		var count int
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if count < offset {
				count++
				continue
			}

			var entity T
			if err := json.Unmarshal(v, &entity); err != nil {
				return fmt.Errorf("error unmarshaling JSON data for key '%s': %v", k, err)
			}

			results = append(results, entity)

			if limit > 0 && len(results) >= limit {
				break
			}
		}

		return nil
	})

	return results, err
}

// getBucketSize retrieves the number of items in a BoltDB bucket using TxStats.
func getBucketSize(db *bolt.DB, bucketName []byte) (int, error) {
	var count int

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found: %s", string(bucketName))
		}

		stats := tx.Bucket(bucketName).Stats()
		count = stats.KeyN

		return nil
	})

	return count, err
}

// addEntity adds an entity to the specified BoltDB bucket with the provided key.
// The entity is first serialized to JSON, and the JSON data is stored in the database.
// If the key already exists, the existing value will be replaced with the new value.
func addEntity[T any](db *bolt.DB, bucket []byte, entity T) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		// Generate ID for the entity using NextSequence
		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		// Use reflection to set the ID field of the entity
		v := reflect.ValueOf(&entity).Elem()
		idField := v.FieldByName("ID")
		if idField.IsValid() && idField.Kind() == reflect.Uint64 {
			idField.SetUint(id)
		} else {
			return errors.New("entity doesn't have a valid ID field of type uint64")
		}

		jsonData, err := json.Marshal(entity)
		if err != nil {
			return err
		}

		err = b.Put(uint64ToBytesBE(id), jsonData)
		return err
	})
}

// updateEntity updates an entity in the specified BoltDB bucket.
// It takes the BoltDB instance, the bucket name, the ID of the entity to update,
// and the entity itself as parameters. The entity is expected to have an "ID" field of type uint64.
// If the entity with the given ID is not found in the bucket, an error is returned.
func updateEntity(db *bolt.DB, bucket []byte, id uint64, entity interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		// Check if the entity with the given ID exists
		existingData := b.Get(uint64ToBytesBE(id))
		if existingData == nil {
			return fmt.Errorf("entity with ID %d not found", id)
		}

		jsonData, err := json.Marshal(entity)
		if err != nil {
			return err
		}

		err = b.Put(uint64ToBytesBE(id), jsonData)
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

// itob returns an 8-byte big endian representation of v.
func uint64ToBytesBE(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	return bytes
}
