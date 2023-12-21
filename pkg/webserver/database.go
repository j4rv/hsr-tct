package webserver

import "github.com/j4rv/hsr-tct/pkg/hsrtct"

var db database

type database interface {
	AddLightCone(lc hsrtct.LightCone) (string, error)
	AddCharacter(c hsrtct.Character) (string, error)
}
