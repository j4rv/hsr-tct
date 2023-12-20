package webserver

import "github.com/j4rv/hsr-tct/pkg/hsrtct"

type database interface {
	Connect() (func(), error)
	AddLightcone(lc hsrtct.LightCone) error
	GetLightcones() ([]hsrtct.LightCone, error)
	AddCharacter(c hsrtct.Character) error
}
