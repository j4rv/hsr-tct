package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/xuri/excelize/v2"
)

const VERSION = "1.1.0"

const FILENAME = "HSRTCT-config.xlsx"
const RESULT_FILENAME = "HSRTCT-results.xlsx"
const LIGHTCONES = "LightCones"
const CHARACTERS = "Characters"
const RELICBUILDS = "RelicBuilds"
const ENEMIES = "Enemies"
const ATTACKS = "Attacks"
const SCENARIOS = "Scenarios"
const EXTERNAL_BUFFS = "ExternalBuffs"
const RESULTS = "HSRTCT Results"

var lightcones map[string]hsrtct.LightCone = map[string]hsrtct.LightCone{}
var characters map[string]hsrtct.Character = map[string]hsrtct.Character{}
var relicbuilds map[string]hsrtct.RelicBuild = map[string]hsrtct.RelicBuild{}
var enemies map[string]hsrtct.Enemy = map[string]hsrtct.Enemy{}
var attacks map[string]hsrtct.Attack = map[string]hsrtct.Attack{}
var externalBuffs map[string][]hsrtct.Buff = map[string][]hsrtct.Buff{}
var scenarios []hsrtct.Scenario = []hsrtct.Scenario{}

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

	log.Println("[INFO] Version: " + VERSION)

	log.Println("[INFO] Reading LightCones...")
	readLightCones(f)
	log.Println("[INFO] Reading Characters...")
	readCharacters(f)
	log.Println("[INFO] Reading RelicBuilds...")
	readRelicBuilds(f)
	log.Println("[INFO] Reading Enemies...")
	readEnemies(f)
	log.Println("[INFO] Reading Attacks...")
	readAttacks(f)
	log.Println("[INFO] Reading External Buffs...")
	readExternalBuffs(f)
	log.Println("[INFO] Reading Scenarios...")
	readScenarios(f)

	log.Println("[INFO] calculating...")
	calcAndWrite()

	fmt.Println("Done!\nPress the Enter key to exit.")
	fmt.Scanln()
}

func calcAndWrite() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet(RESULTS)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)
	if err != nil {
		fmt.Println(err)
		return
	}

	centeredNumberStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 2,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellValue(RESULTS, "A1", "Scenario")
	f.SetCellValue(RESULTS, "B1", "Damage")
	f.SetColWidth(RESULTS, "A", "A", 80)
	f.SetColWidth(RESULTS, "B", "B", 20)
	f.SetColStyle(RESULTS, "B", centeredNumberStyle)

	for rowIndex, scenario := range scenarios {
		rowIndex++
		f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 0), scenario.Name)

		dmg, err := hsrtct.CalcAvgDmgScenario(scenario)
		if err != nil {
			log.Println("[ERROR] failed to calculate damage for scenario: " + scenario.Name + ", " + err.Error())
			f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 1), "failed to calculate damage for scenario: "+scenario.Name+", "+err.Error())
		} else {
			log.Println("[INFO] " + scenario.Name + ": " + strconv.FormatFloat(dmg, 'f', 2, 64))
			f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 1), dmg)
		}

		explanation, err := hsrtct.ExplainDmgScenario(scenario)
		if err != nil {
			log.Println("[ERROR] failed to explain damage for scenario: " + scenario.Name + ", " + err.Error())
			explanation = "failed to explain damage for scenario: " + scenario.Name + ", " + err.Error()
		}
		stats, err := hsrtct.ExplainFinalStats(scenario.Character, scenario.LightCone, scenario.RelicBuild)
		if err != nil {
			log.Println("[ERROR] failed to explain final stats for scenario: " + scenario.Name + ", " + err.Error())
			explanation = "failed to explain final stats for scenario: " + scenario.Name + ", " + err.Error()
		}
		f.AddComment(RESULTS, excelize.Comment{
			Author: "HSRTCT",
			Cell:   spreadsheetCoordinate(rowIndex, 1),
			Text:   fmt.Sprintf("Stats:\n%s\n\n%s", stats, explanation),
		})
	}

	if err := f.SaveAs(RESULT_FILENAME); err != nil {
		log.Println("[ERROR] failed to save results: " + err.Error())
		fmt.Println("failed to save results: " + err.Error())
	}
}

func spreadsheetCoordinate(row, col int) string {
	columnLetters := ""
	col++
	for col > 0 {
		col--
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
