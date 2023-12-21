package bolt

import (
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

var charactersBucket = []byte("characters")

func (database *database) AddCharacter(c hsrtct.Character) error {
	return addEntity(database.db, charactersBucket, c)
}

func (database *database) GetCharacter(id string) (hsrtct.Character, error) {
	return getEntity[hsrtct.Character](database.db, charactersBucket, id)
}

func (database *database) DeleteCharacter(id string) error {
	return deleteEntity(database.db, charactersBucket, id)
}
