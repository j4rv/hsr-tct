package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/xuri/excelize/v2"
)

const FILENAME = "HSRTCT-config.xlsx"
const LIGHTCONES = "LightCones"
const CHARACTERS = "Characters"
const RELICBUILDS = "RelicBuilds"

var lightcones map[string]hsrtct.LightCone = map[string]hsrtct.LightCone{}
var characters map[string]hsrtct.Character = map[string]hsrtct.Character{}
var relicbuilds map[string]hsrtct.RelicBuild = map[string]hsrtct.RelicBuild{}

func main() {
	f, err := excelize.OpenFile(FILENAME)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	readLightCones(f)
	readCharacters(f)
	readRelicBuilds(f)
}

func readLightCones(f *excelize.File) {
	log.Println("readLightCones...")
	rows, err := f.GetRows(LIGHTCONES)
	if err != nil {
		panic("failed to read LightCones: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		lc := hsrtct.LightCone{}
		lc.Name = row[0]
		lc.Level = 80
		lc.BaseHp = mustParseFloat(row[1])
		lc.BaseAtk = mustParseFloat(row[2])
		lc.BaseDef = mustParseFloat(row[3])
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
	log.Println("readCharacters...")
	rows, err := f.GetRows(CHARACTERS)
	if err != nil {
		panic("failed to read Characters: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		c := hsrtct.Character{}
		c.Name = row[0]
		c.Level = mustParseInt(row[1])
		c.BaseHp = mustParseFloat(row[2])
		c.BaseAtk = mustParseFloat(row[3])
		c.BaseDef = mustParseFloat(row[4])
		c.BaseSpd = mustParseFloat(row[5])
		c.BaseAggro = mustParseFloat(row[6])
		c.Element = hsrtct.Element(row[7])
		for j := 0; j < 4; j++ {
			b, err := readBuff(f, CHARACTERS, i, 8+j*4)
			if err == nil {
				c.Buffs = append(c.Buffs, b)
			}
		}
		characters[c.Name] = c
	}
}

func readRelicBuilds(f *excelize.File) {
	log.Println("readRelicBuilds...")
	rows, err := f.GetRows(RELICBUILDS)
	if err != nil {
		panic("failed to read RelicBuilds: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 || row[0] == "" {
			continue
		}
		rb := hsrtct.RelicBuild{}
		rb.Name = row[0]
		rb.Relics[0] = hsrtct.Relic{MainStat: hsrtct.Stat(row[1])}
		rb.Relics[1] = hsrtct.Relic{MainStat: hsrtct.Stat(row[2])}
		rb.Relics[2] = hsrtct.Relic{MainStat: hsrtct.Stat(row[3])}
		rb.Relics[3] = hsrtct.Relic{MainStat: hsrtct.Stat(row[4])}
		rb.Relics[4] = hsrtct.Relic{MainStat: hsrtct.Stat(row[5])}
		rb.Relics[5] = hsrtct.Relic{MainStat: hsrtct.Stat(row[6])}

		// columns 7 to 18 are substat roll amounts, column 19 is roll type
		for j := 7; j < 19; j++ {
			substat := hsrtct.RelicSubstat{}
			substat.RollType = hsrtct.RollType(row[19])
			substat.Stat = hsrtct.Stat(rows[0][j])
			substat.Rolls = mustParseInt(row[j])
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

func readBuff(f *excelize.File, sheetName string, row, col int) (hsrtct.Buff, error) {
	rawValue, err := f.GetCellValue(sheetName, spreadsheetCoordinate(row, col))
	if err != nil {
		return hsrtct.Buff{}, err
	}
	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
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

func spreadsheetCoordinate(row, col int) string {
	columnLetters := ""
	for col > 0 {
		columnLetters = string('A'+col%26) + columnLetters
		col /= 26
	}
	return fmt.Sprintf("%s%d", columnLetters, row+1)
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
func mustParseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
