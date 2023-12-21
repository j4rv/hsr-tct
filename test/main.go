package main

import (
	"github.com/j4rv/hsr-tct/pkg/bolt"
	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

func main() {
	db := bolt.New()
	close, err := db.Init("my.db")
	if err != nil {
		panic(err)
	}
	defer close()

	lc := hsrtct.LightCone{
		Name: "LightCone",
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 26},
			{Stat: hsrtct.DmgBonus, Value: 20},
		},
	}
	err = db.AddLightCone(lc)
	lc = hsrtct.LightCone{
		Name: "LightCone 2",
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 26},
			{Stat: hsrtct.CritRate, Value: 16},
		},
	}
	err = db.AddLightCone(lc)
	if err != nil {
		panic(err)
	}

}
