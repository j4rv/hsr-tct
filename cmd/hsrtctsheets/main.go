package main

import (
	"fmt"
	"strconv"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/xuri/excelize/v2"
)

var lightcones map[string]hsrtct.LightCone = map[string]hsrtct.LightCone{}

func main() {
	f, err := excelize.OpenFile("HSRTCT-config.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	readLightCones(f)
}

func readLightCones(f *excelize.File) {
	rows, err := f.GetRows("LightCones")
	if err != nil {
		panic("failed to read LightCones: " + err.Error())
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		lc := hsrtct.LightCone{}
		lc.Name = row[0]
		lc.Level = 80
		lc.BaseHp = mustParseFloat(row[1])
		lc.BaseAtk = mustParseFloat(row[2])
		lc.BaseDef = mustParseFloat(row[3])
		for j := 0; j < 4; j++ {
			b, err := readBuff(f, "LightCones", i, 5+j*4)
			if err == nil {
				lc.Buffs = append(lc.Buffs, b)
			}
		}
		lightcones[lc.Name] = lc
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
		col -= 1
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
