package webserver

import "github.com/j4rv/hsr-tct/pkg/hsrtct"

var db database

type database interface {
	AddLightCone(lc hsrtct.LightCone) error
	GetLightCones() ([]hsrtct.LightCone, error)
}
