package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

var lightconesBucket = []byte("lightcones")

func (database *database) AddLightCone(lc hsrtct.LightCone) error {
	return addEntity(database.db, lightconesBucket, lc)
}

func (database *database) GetLightCone(id string) (hsrtct.LightCone, error) {
	return getEntity[hsrtct.LightCone](database.db, lightconesBucket, id)
}

func (database *database) GetLightCones() ([]hsrtct.LightCone, error) {
	return getAllEntities[hsrtct.LightCone](database.db, lightconesBucket)
}

func (database *database) DeleteLightcone(id string) error {
	return deleteEntity(database.db, lightconesBucket, id)
}
