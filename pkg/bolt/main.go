package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	uuid "github.com/satori/go.uuid"
)

type database struct {
	db *bolt.DB
}

func New() *database {
	return &database{}
}

var charactersBucket = []byte("Characters")
var lightconesBucket = []byte("Lightcones")
var enemiesBucket = []byte("Enemies")

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

func (database *database) AddLightCone(c hsrtct.LightCone) (string, error) {
	c.ID = uuid.NewV4().String()
	return c.ID, addEntity(database.db, lightconesBucket, c.ID, c)
}

func (database *database) AddCharacter(c hsrtct.Character) (string, error) {
	c.ID = uuid.NewV4().String()
	return c.ID, addEntity(database.db, charactersBucket, c.ID, c)
}

func (database *database) GetCharacter(id string) (hsrtct.Character, error) {
	return getEntity[hsrtct.Character](database.db, charactersBucket, id)
}

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

func NewUUID() []byte {
	return uuid.NewV4().Bytes()
}
