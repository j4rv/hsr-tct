package main

import (
	"log"

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

	id, err := db.AddCharacter(hsrtct.Character{
		Name:    "Hook",
		Level:   1,
		Element: hsrtct.Fire,
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 4 + 6 + 6 + 8},
			{Stat: hsrtct.HpPct, Value: 4 + 6 + 8},
			{Stat: hsrtct.CritDmg, Value: 5.3 + 8},
			{Stat: hsrtct.DmgBonus, Value: 20},
			{Stat: hsrtct.DmgBonus, Value: 20, DamageTag: hsrtct.Skill},
		},
		LightCone: lc,
	})
	if err != nil {
		panic(err)
	}

	c, err := db.GetCharacter(id)
	if err != nil {
		panic(err)
	}

	log.Println(c)
}
