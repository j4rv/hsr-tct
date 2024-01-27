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
const ENEMIES = "Enemies"
const ATTACKS = "Attacks"
const SCENARIOS = "Scenarios"

var lightcones map[string]hsrtct.LightCone = map[string]hsrtct.LightCone{}
var characters map[string]hsrtct.Character = map[string]hsrtct.Character{}
var relicbuilds map[string]hsrtct.RelicBuild = map[string]hsrtct.RelicBuild{}
var enemies map[string]hsrtct.Enemy = map[string]hsrtct.Enemy{}
var attacks map[string]hsrtct.Attack = map[string]hsrtct.Attack{}
var scenarios map[string]hsrtct.Scenario = map[string]hsrtct.Scenario{}

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

	log.Println("readLightCones...")
	readLightCones(f)
	log.Println("readCharacters...")
	readCharacters(f)
	log.Println("readRelicBuilds...")
	readRelicBuilds(f)
	log.Println("readEnemies...")
	readEnemies(f)
	log.Println("readAttacks...")
	readAttacks(f)
	log.Println("readScenarios...")
	readScenarios(f)

	log.Println(scenarios["Hook Prisoner Solo 2dots"].Attacks)
	log.Println(scenarios["Hook Prisoner Solo 2dots"].Enemies)
	log.Println(scenarios["Hook Prisoner Solo 2dots"].FocusedEnemy)

	for k, v := range scenarios {
		dmg, err := hsrtct.CalcAvgDmgScenario(v)
		if err != nil {
			panic("failed to calculate damage for scenario: " + k + ", " + err.Error())
		}
		fmt.Println(k, dmg)
		explanation, err := hsrtct.ExplainDmgScenario(v)
		if err != nil {
			panic("failed to explain damage for scenario: " + k + ", " + err.Error())
		}
		fmt.Println(explanation)
	}
}

func spreadsheetCoordinate(row, col int) string {
	columnLetters := ""
	for col > 0 {
		columnLetters = string(rune('A'+col%26)) + columnLetters
		col /= 26
	}
	return fmt.Sprintf("%s%d", columnLetters, row+1)
}

func mustParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
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
