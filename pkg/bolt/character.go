package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	uuid "github.com/satori/go.uuid"
)

var charactersBucket = []byte("characters")

func (database *database) AddCharacter(c hsrtct.Character) (string, error) {
	c.ID = uuid.NewV4().String()
	return c.ID, addEntity(database.db, charactersBucket, c.ID, c)
}

func (database *database) GetCharacter(id string) (hsrtct.Character, error) {
	return getEntity[hsrtct.Character](database.db, charactersBucket, id)
}

func (database *database) UpdateCharacter(id string, c hsrtct.Character) error {
	return addEntity(database.db, charactersBucket, id, c)
}

func (database *database) DeleteCharacter(id string) error {
	return deleteEntity(database.db, charactersBucket, id)
}
