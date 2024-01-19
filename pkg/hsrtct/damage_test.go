package hsrtct_test

import (
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
	hook.LightCone = GetAeonLC()
	hook.RelicBuild = GetHookRelicBuild()
	enemy := GetBasicEnemy()
	attack := hsrtct.Attack{
		ScalingStat: hsrtct.Atk,
		Multiplier:  432 + 110,
		Element:     hsrtct.Fire,
		AttackTag:   hsrtct.Ultimate,
	}
	damage, err := hsrtct.CalcAvgDamageST(hook, enemy, attack)
	assertNotNil(t, err)
	if int(damage) != 41425 {
		t.Fatalf("Expected damage to be 41425, got %v", damage)
	}
}

func TestCalcAvgDamageMultipleEnemies(t *testing.T) {
	hook := GetHookCharacter()
	hook.LightCone = GetAeonLC()
	hook.RelicBuild = GetHookRelicBuild()
	leftEnemy := GetBasicEnemy()
	centerEnemy := GetBasicEnemy()
	centerEnemy.Buffs = append(centerEnemy.Buffs, hsrtct.Buff{Stat: hsrtct.DefShred, Value: 45 + 8})
	rightEnemy := GetBasicEnemy()

	skillCenter := hsrtct.Attack{
		ScalingStat: hsrtct.Atk,
		Multiplier:  308 + 110,
		Element:     hsrtct.Fire,
		AttackTag:   hsrtct.Skill,
	}
	skillSplash := hsrtct.Attack{
		ScalingStat: hsrtct.Atk,
		Multiplier:  88 + 110,
		Element:     hsrtct.Fire,
		AttackTag:   hsrtct.Skill,
	}

	leftDmg, err := hsrtct.CalcAvgDamageST(hook, leftEnemy, skillSplash)
	assertNotNil(t, err)
	rightDmg, err := hsrtct.CalcAvgDamageST(hook, rightEnemy, skillSplash)
	assertNotNil(t, err)
	centerDmg, err := hsrtct.CalcAvgDamageST(hook, centerEnemy, skillCenter)
	assertNotNil(t, err)

	damage := leftDmg + centerDmg + rightDmg
	if int(damage) != 91656 {
		t.Fatalf("Expected damage to be 91656, got %v", damage)
	}
}

func assertNotNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected non-nil error, got nil")
	}
}
