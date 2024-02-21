package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/xuri/excelize/v2"
)

const VERSION = "1.2.0"

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
	f.SetCellValue(RESULTS, "C1", "Explanation page name")
	f.SetColWidth(RESULTS, "A", "A", 150)
	f.SetColWidth(RESULTS, "B", "B", 20)
	f.SetColWidth(RESULTS, "C", "C", 20)
	f.SetColStyle(RESULTS, "B", centeredNumberStyle)

	for rowIndex, scenario := range scenarios {
		rowIndex++
		explanationSheetName := fmt.Sprintf("SCN %d", rowIndex)
		f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 0), scenario.Name)
		f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 2), explanationSheetName)

		result, err := hsrtct.CalcAvgDmgScenario(scenario)

		if err != nil {
			log.Println("[ERROR] failed to calculate damage for scenario: " + scenario.Name + ", " + err.Error())
			f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 1), "Failed to calculate damage for scenario: "+scenario.Name+", "+err.Error())
		} else {
			formattedDmg := strconv.FormatFloat(result.TotalDmg, 'f', 0, 64)
			log.Println("[INFO] " + scenario.Name + ": " + formattedDmg)
			f.SetCellValue(RESULTS, spreadsheetCoordinate(rowIndex, 1), formattedDmg)

			for expIndex, exp := range result.Explanations {
				f.NewSheet(explanationSheetName)
				f.SetColWidth(explanationSheetName, "A", "Z", 40)
				for i, expLine := range strings.Split(exp, "\n") {
					f.SetCellValue(explanationSheetName, spreadsheetCoordinate(i, expIndex), expLine)
				}
			}
		}
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
