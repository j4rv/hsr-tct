package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

var enemiesBucket = []byte("enemies")

func (database *database) GetEnemy(id string) (hsrtct.Enemy, error) {
	return getEntity[hsrtct.Enemy](database.db, enemiesBucket, id)
}

func (database *database) DeleteEnemy(id string) error {
	return deleteEntity(database.db, enemiesBucket, id)
}

func (database *database) AddEnemy(e hsrtct.Enemy) error {
	return addEntity(database.db, enemiesBucket, e)
}
