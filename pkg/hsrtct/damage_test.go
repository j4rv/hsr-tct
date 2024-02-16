package hsrtct_test

import (
	"log"
	"testing"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

func GetHookCharacter() hsrtct.Character {
	return hsrtct.Character{
		Name:      "Hook",
		Level:     80,
		BaseHp:    1340,
		BaseAtk:   617,
		BaseDef:   352,
		BaseSpd:   94,
		BaseAggro: 125,
		Element:   hsrtct.Fire,
		// Kit, traces and eidolons
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 4 + 6 + 6 + 8},
			{Stat: hsrtct.HpPct, Value: 4 + 6 + 8},
			{Stat: hsrtct.CritDmg, Value: 5.3 + 8},
			{Stat: hsrtct.DmgBonus, Value: 20},
			{Stat: hsrtct.DmgBonus, Value: 20, DamageTag: hsrtct.Skill}, // FIXME: Add buffs to Attack as well, this only applies to her enhanced skill
		},
	}
}

func GetAeonLC() hsrtct.LightCone {
	return hsrtct.LightCone{
		Name:    "On the Fall of an Aeon",
		Level:   80,
		BaseHp:  1058,
		BaseAtk: 529,
		BaseDef: 396,
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 16 * 4},
			{Stat: hsrtct.DmgBonus, Value: 24},
		},
	}
}

func GetHookRelicBuild() hsrtct.RelicBuild {
	build := hsrtct.RelicBuild{}
	rollType := hsrtct.RollTypeAvg

	build.Relics[0] = hsrtct.Relic{MainStat: hsrtct.Hp}
	build.Relics[1] = hsrtct.Relic{MainStat: hsrtct.Atk}
	build.Relics[2] = hsrtct.Relic{MainStat: hsrtct.CritRate}
	build.Relics[3] = hsrtct.Relic{MainStat: hsrtct.AtkPct}
	build.Relics[4] = hsrtct.Relic{MainStat: hsrtct.DmgBonus}
	build.Relics[5] = hsrtct.Relic{MainStat: hsrtct.AtkPct}

	build.SubStats = []hsrtct.RelicSubstat{
		{Stat: hsrtct.Atk, Rolls: 2, RollType: rollType},
		{Stat: hsrtct.AtkPct, Rolls: 8, RollType: rollType},
		{Stat: hsrtct.CritRate, Rolls: 10, RollType: rollType},
		{Stat: hsrtct.CritDmg, Rolls: 12, RollType: rollType},
	}

	build.SetEffects = []hsrtct.Buff{
		{Stat: hsrtct.AtkPct, Value: 12},
		{Stat: hsrtct.DefIgnore, Value: 6 * 2},
		{Stat: hsrtct.CritRate, Value: 8},
		{Stat: hsrtct.DmgBonus, Value: 20, DamageTag: hsrtct.Basic},
		{Stat: hsrtct.DmgBonus, Value: 20, DamageTag: hsrtct.Skill},
	}

	return build
}

func GetBasicEnemy() hsrtct.Enemy {
	return hsrtct.Enemy{
		Name:  "Basic",
		Level: 85,
		Buffs: []hsrtct.Buff{},
	}
}

func TestCalcAvgDamageUltimate(t *testing.T) {
	hook := GetHookCharacter()
	lc := GetAeonLC()
	rb := GetHookRelicBuild()
	enemy := GetBasicEnemy()
	attack := hsrtct.Attack{
		ScalingStat: hsrtct.Atk,
		Multiplier:  432 + 110,
		Element:     hsrtct.Fire,
		DamageTag:   hsrtct.Ultimate,
	}

	scn := hsrtct.Scenario{
		Character:    hook,
		LightCone:    lc,
		RelicBuild:   rb,
		Enemies:      []hsrtct.Enemy{enemy},
		Attacks:      map[*hsrtct.Attack]float64{&attack: 1},
		FocusedEnemy: 0,
	}

	scnResult, err := hsrtct.CalcAvgDmgScenario(scn)
	assertNilError(t, err)
	if int(scnResult.TotalDmg) != 41425 {
		t.Fatalf("Expected damage to be 41425, got %v", scnResult.TotalDmg)
	}
}

func TestCalcAvgDamageMultipleEnemies(t *testing.T) {
	hook := GetHookCharacter()
	lc := GetAeonLC()
	rb := GetHookRelicBuild()
	leftEnemy := GetBasicEnemy()
	centerEnemy := GetBasicEnemy()
	centerEnemy.Buffs = append(centerEnemy.Buffs, hsrtct.Buff{Stat: hsrtct.DefShred, Value: 45 + 8})
	rightEnemy := GetBasicEnemy()

	skill := hsrtct.Attack{
		ScalingStat:      hsrtct.Atk,
		AttackAOE:        hsrtct.Blast,
		Multiplier:       308 + 110,
		MultiplierSplash: 88 + 110,
		Element:          hsrtct.Fire,
		DamageTag:        hsrtct.Skill,
	}

	scn := hsrtct.Scenario{
		Character:    hook,
		LightCone:    lc,
		RelicBuild:   rb,
		Enemies:      []hsrtct.Enemy{leftEnemy, centerEnemy, rightEnemy},
		Attacks:      map[*hsrtct.Attack]float64{&skill: 1},
		FocusedEnemy: 1,
	}

	scnResult, err := hsrtct.CalcAvgDmgScenario(scn)
	assertNilError(t, err)

	damage := scnResult.TotalDmg
	if int(damage) != 91656 {
		t.Fatalf("Expected damage to be 91656, got %v", damage)
	}
}

func TestExplainFinalStats(t *testing.T) {
	hook := GetHookCharacter()
	lc := GetAeonLC()
	rb := GetHookRelicBuild()
	log.Println(hsrtct.FinalStats(hook, lc, rb, GetBasicEnemy(), hsrtct.Attack{}))
}

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected nil error, got '%v'", err)
	}
}
