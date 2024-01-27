package main

import (
	"log"
	"strconv"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/xuri/excelize/v2"
)

func readLightCones(f *excelize.File) {
	rows, err := f.GetRows(LIGHTCONES)
	if err != nil {
		panic("failed to read LightCones: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		lc := hsrtct.LightCone{
			Name:    row[0],
			Level:   80,
			BaseHp:  mustParseFloat(row[1]),
			BaseAtk: mustParseFloat(row[2]),
			BaseDef: mustParseFloat(row[3]),
		}
		for j := 0; j < 4; j++ {
			b, err := readBuff(f, LIGHTCONES, i, 4+j*4)
			if err == nil {
				lc.Buffs = append(lc.Buffs, b)
			}
		}
		lightcones[lc.Name] = lc
	}
}

func readCharacters(f *excelize.File) {
	rows, err := f.GetRows(CHARACTERS)
	if err != nil {
		panic("failed to read Characters: " + err.Error())
	}

	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}

		character := hsrtct.Character{
			Name:      row[0],
			Level:     mustParseInt(row[1]),
			BaseHp:    mustParseFloat(row[2]),
			BaseAtk:   mustParseFloat(row[3]),
			BaseDef:   mustParseFloat(row[4]),
			BaseSpd:   mustParseFloat(row[5]),
			BaseAggro: mustParseFloat(row[6]),
			Element:   hsrtct.Element(row[7]),
		}

		for j := 0; j < 4; j++ {
			buff, err := readBuff(f, CHARACTERS, i, 8+j*4)
			if err == nil {
				character.Buffs = append(character.Buffs, buff)
			}
		}

		characters[character.Name] = character
	}
}

func readRelicBuilds(f *excelize.File) {
	rows, err := f.GetRows(RELICBUILDS)
	if err != nil {
		panic("failed to read RelicBuilds: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		rb := hsrtct.RelicBuild{
			Name: row[0],
			Relics: [6]hsrtct.Relic{
				{MainStat: hsrtct.Stat(row[1])},
				{MainStat: hsrtct.Stat(row[2])},
				{MainStat: hsrtct.Stat(row[3])},
				{MainStat: hsrtct.Stat(row[4])},
				{MainStat: hsrtct.Stat(row[5])},
				{MainStat: hsrtct.Stat(row[6])},
			},
			SubStats:   make([]hsrtct.RelicSubstat, 0, 12),
			SetEffects: make([]hsrtct.Buff, 0),
		}

		for j := 7; j < 19; j++ {
			substat := hsrtct.RelicSubstat{
				RollType: hsrtct.RollType(row[19]),
				Stat:     hsrtct.Stat(rows[0][j]),
				Rolls:    mustParseInt(row[j]),
			}
			rb.SubStats = append(rb.SubStats, substat)
		}

		for j := 0; j < 6; j++ {
			b, err := readBuff(f, RELICBUILDS, i, 20+j*4)
			if err == nil {
				rb.SetEffects = append(rb.SetEffects, b)
			}
		}

		relicbuilds[rb.Name] = rb
	}
}

func readEnemies(f *excelize.File) {
	rows, err := f.GetRows(ENEMIES)
	if err != nil {
		panic("failed to read Enemies: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		enemy := hsrtct.Enemy{
			Name:  row[0],
			Level: mustParseInt(row[1]),
		}

		for j := 0; j < 5; j++ {
			buff, err := readBuff(f, ENEMIES, i, 2+j*4)
			if err == nil {
				enemy.Buffs = append(enemy.Buffs, buff)
			}
		}

		enemies[enemy.Name] = enemy
	}
}

func readAttacks(f *excelize.File) {
	rows, err := f.GetRows(ATTACKS)
	if err != nil {
		panic("failed to read Attacks: " + err.Error())
	}

	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}

		aoe, err := hsrtct.ParseAttackAOE(row[6])
		if err != nil {
			panic("failed to read Attacks: " + err.Error())
		}
		attack := hsrtct.Attack{
			Name:             row[0],
			ScalingStat:      hsrtct.Stat(row[1]),
			Multiplier:       mustParseFloat(row[2]),
			MultiplierSplash: mustParseFloat(row[3]),
			Element:          hsrtct.Element(row[4]),
			DamageTag:        hsrtct.DamageTag(row[5]),
			AttackAOE:        aoe,
		}

		for j := 0; j < 5; j++ {
			buff, err := readBuff(f, ATTACKS, i, 7+j*4)
			if err == nil {
				attack.Buffs = append(attack.Buffs, buff)
			}
		}

		attacks[attack.Name] = attack
	}
}

func readScenarios(f *excelize.File) {
	rows, err := f.GetRows(SCENARIOS)
	if err != nil {
		panic("failed to read Scenarios: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" || row[2] != "TRUE" {
			continue
		}

		scenario := hsrtct.Scenario{
			Name:         row[0],
			Notes:        row[1],
			Character:    characters[row[3]],
			LightCone:    lightcones[row[4]],
			RelicBuild:   relicbuilds[row[5]],
			Enemies:      make([]hsrtct.Enemy, 0),
			FocusedEnemy: mustParseInt(row[11]) - 1, // Index in excel starts at 1
			Attacks:      make(map[*hsrtct.Attack]float64),
		}

		for j := 0; j < 5; j++ {
			enemyName := row[6+j]
			if enemyName == "" {
				continue
			}
			enemy := enemies[enemyName]
			scenario.Enemies = append(scenario.Enemies, enemy)
		}

		for j := 0; j < 8; j++ {
			attack, mult := readAttack(f, SCENARIOS, i, 12+j*2)
			if attack == nil || mult == 0 {
				continue
			}
			scenario.Attacks[attack] = mult
		}

		scenarios[scenario.Name] = scenario
	}
}

func readBuff(f *excelize.File, sheetName string, row, col int) (hsrtct.Buff, error) {
	rawValue, err := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col))
	if rawValue == "" || err != nil {
		return hsrtct.Buff{}, err
	}
	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		log.Println("Couldnt parse as float: " + rawValue)
		return hsrtct.Buff{}, err
	}
	stat, _ := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col+1))
	dmgTag, _ := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col+2))
	element, _ := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col+3))
	return hsrtct.Buff{
		Value:     value,
		Stat:      hsrtct.Stat(stat),
		DamageTag: hsrtct.DamageTag(dmgTag),
		Element:   hsrtct.Element(element),
	}, nil
}

func readAttack(f *excelize.File, sheetName string, row, col int) (*hsrtct.Attack, float64) {
	attackName, err := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col))
	if attackName == "" || err != nil {
		return nil, 0
	}
	attack := attacks[attackName]
	rawMult, err := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col+1))
	if rawMult == "" || err != nil {
		return nil, 0
	}
	return &attack, mustParseFloat(rawMult)
}
