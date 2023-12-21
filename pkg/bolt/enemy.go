package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	uuid "github.com/satori/go.uuid"
)

var enemiesBucket = []byte("enemies")

func (database *database) GetEnemy(id string) (hsrtct.Enemy, error) {
	return getEntity[hsrtct.Enemy](database.db, enemiesBucket, id)
}

func (database *database) UpdateEnemy(id string, e hsrtct.Enemy) error {
	return addEntity(database.db, enemiesBucket, id, e)
}

func (database *database) DeleteEnemy(id string) error {
	return deleteEntity(database.db, enemiesBucket, id)
}

func (database *database) AddEnemy(e hsrtct.Enemy) (string, error) {
	e.ID = uuid.NewV4().String()
	return e.ID, addEntity(database.db, enemiesBucket, e.ID, e)
}
