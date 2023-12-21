package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	uuid "github.com/satori/go.uuid"
)

var lightconesBucket = []byte("lightcones")

func (database *database) AddLightCone(lc hsrtct.LightCone) (string, error) {
	lc.ID = uuid.NewV4().String()
	return lc.ID, addEntity(database.db, lightconesBucket, lc.ID, lc)
}

func (database *database) GetLightCone(id string) (hsrtct.LightCone, error) {
	return getEntity[hsrtct.LightCone](database.db, lightconesBucket, id)
}

func (database *database) UpdateLightCone(id string, lc hsrtct.LightCone) error {
	return addEntity(database.db, lightconesBucket, id, lc)
}

func (database *database) DeleteLightcone(id string) error {
	return deleteEntity(database.db, lightconesBucket, id)
}
