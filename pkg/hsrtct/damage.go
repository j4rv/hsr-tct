package hsrtct

import (
	"errors"
	"fmt"
	"math"
)

var ErrInvalidScalingStat = errors.New("invalid scaling stat")

type DamageTag string

const (
	AnyAttack DamageTag = ""
	Basic     DamageTag = "Basic"
	Skill     DamageTag = "Skill"
	Ultimate  DamageTag = "Ultimate"
	FollowUp  DamageTag = "FollowUp"
	Dot       DamageTag = "Dot"
)

func (d DamageTag) Is(tag DamageTag) bool {
	return d == tag || d == AnyAttack || tag == AnyAttack
}

func AttackTagKeys() []DamageTag {
	return []DamageTag{Basic, Skill, Ultimate, FollowUp, Dot}
}

type AttackAOE string

const (
	Single AttackAOE = ""
	Blast  AttackAOE = "Blast"
	All    AttackAOE = "All"
)

// Attack is an attack that can be used by a character.
// If AOE is Single: will use Multiplier and MultiplierSplash will be ignored.
// If AOE is Blast, will use Multiplier for the focused enemy and MultiplierSplash for its neighbors.
// If AOE is All, will use Multiplier for all enemies.
type Attack struct {
	ID               uint64
	Name             string
	ScalingStat      Stat
	Multiplier       float64
	MultiplierSplash float64
	Element          Element
	AttackTag        DamageTag
	AttackAOE        AttackAOE
	Buffs            []Buff
}

type Scenario struct {
	ID           uint64
	Name         string
	Notes        string
	Character    Character
	Enemies      []Enemy
	FocusedEnemy int
	Attacks      []Attack
}

func CalcAvgDmgScenario(s Scenario) (float64, error) {
	totalDmg := 0.0
	for _, attack := range s.Attacks {
		switch attack.AttackAOE {

		case Single:
			dmg, err := CalcAvgDamage(s.Character, s.Enemies[s.FocusedEnemy], attack, false)
			if err != nil {
				return 0, err
			}
			totalDmg += dmg

		case Blast:
			avgDmg, err := CalcAvgDamage(s.Character, s.Enemies[s.FocusedEnemy], attack, false)
			if err != nil {
				return 0, err
			}
			totalDmg += avgDmg
			if s.FocusedEnemy-1 >= 0 {
				splashDmg, err := CalcAvgDamage(s.Character, s.Enemies[s.FocusedEnemy-1], attack, true)
				if err != nil {
					return 0, err
				}
				totalDmg += splashDmg
			}
			if s.FocusedEnemy+1 < len(s.Enemies) {
				splashDmg, err := CalcAvgDamage(s.Character, s.Enemies[s.FocusedEnemy+1], attack, true)
				if err != nil {
					return 0, err
				}
				totalDmg += splashDmg
			}

		case All:
			for _, enemy := range s.Enemies {
				dmg, err := CalcAvgDamage(s.Character, enemy, attack, false)
				if err != nil {
					return 0, err
				}
				totalDmg += dmg
			}
		}
	}
	return totalDmg, nil
}

func CalcAvgDamage(c Character, e Enemy, a Attack, isSplash bool) (float64, error) {
	baseDamage, err := CalcBaseDamage(c, e, a, isSplash)
	if err != nil {
		return 0, err
	}
	critMult := CalcAvgCritMultiplier(c, e, a)
	dmgBonusMult := CalcDmgBonusMult(c, e, a)
	resMult := CalcResistanceMultiplier(c, e, a)
	defMult := CalcDefenseMultiplier(c, e, a)
	vulnMult := CalcVulnerabilityMultiplier(c, e, a)
	return baseDamage * critMult * dmgBonusMult * resMult * defMult * vulnMult, nil
}

func ExplainDamage(c Character, e Enemy, a Attack, isSplash bool) (string, error) {
	baseDamage, err := CalcBaseDamage(c, e, a, isSplash)
	if err != nil {
		return "", err
	}
	critMult := CalcAvgCritMultiplier(c, e, a)
	dmgBonusMult := CalcDmgBonusMult(c, e, a)
	resMult := CalcResistanceMultiplier(c, e, a)
	defMult := CalcDefenseMultiplier(c, e, a)
	vulnMult := CalcVulnerabilityMultiplier(c, e, a)
	return fmt.Sprintf(
		"Base Damage: %.2f\n"+
			"Crit Multiplier: %.2f\n"+
			"Damage Bonus Multiplier: %.2f\n"+
			"Resistance Multiplier: %.2f\n"+
			"Defense Multiplier: %.2f\n"+
			"Vulnerability Multiplier: %.2f",
		baseDamage, critMult, dmgBonusMult, resMult, defMult, vulnMult), nil
}

func CalcBaseDamage(c Character, e Enemy, a Attack, isSplash bool) (float64, error) {
	baseDamage := 0.0
	mult := a.Multiplier
	if isSplash {
		mult = a.MultiplierSplash
	}
	switch a.ScalingStat {
	case Hp:
		baseDamage = c.FinalStatValue(Hp, a.AttackTag, a.Element, a.Buffs) * mult / 100
	case Atk:
		baseDamage = c.FinalStatValue(Atk, a.AttackTag, a.Element, a.Buffs) * mult / 100
	case Def:
		baseDamage = c.FinalStatValue(Def, a.AttackTag, a.Element, a.Buffs) * mult / 100
	default:
		return 0, ErrInvalidScalingStat
	}
	return baseDamage, nil
}

func CalcAvgCritMultiplier(c Character, e Enemy, a Attack) float64 {
	if a.AttackTag == Dot {
		return 1
	}
	critRate := c.FinalStatValue(CritRate, a.AttackTag, a.Element, a.Buffs)
	critDamage := c.FinalStatValue(CritDmg, a.AttackTag, a.Element, a.Buffs)
	critRate = math.Min(critRate, 100)
	return 1 + (critRate / 100 * critDamage / 100)
}

func CalcDmgBonusMult(c Character, e Enemy, a Attack) float64 {
	return 1 + c.FinalStatValue(DmgBonus, a.AttackTag, a.Element, a.Buffs)/100
}

func CalcResistanceMultiplier(c Character, e Enemy, a Attack) float64 {
	res := 0.0

	for _, buff := range e.Buffs {
		if buff.Stat == ElementalRes && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			res += buff.Value
		}
		if buff.Stat == ResShred && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			res -= buff.Value
		}
	}

	for _, buff := range c.AllBuffs() {
		if buff.Stat == ResPen && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			res -= buff.Value
		}
	}

	for _, buff := range a.Buffs {
		if buff.Stat == ResPen && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			res -= buff.Value
		}
	}

	return 1.0 - res/100
}

func CalcDefenseMultiplier(c Character, e Enemy, a Attack) float64 {
	flatDef := 0.0
	defPct := 0.0
	defReduction := 0.0
	baseDef := 200.0 + 10.0*float64(e.Level)

	for _, buff := range e.Buffs {
		if buff.Stat == Def {
			flatDef += buff.Value
		}
		if buff.Stat == DefPct {
			defPct += buff.Value
		}
		if buff.Stat == DefShred && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			defReduction += buff.Value
		}
	}

	for _, buff := range c.AllBuffs() {
		if buff.Stat == DefIgnore && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			defReduction += buff.Value
		}
	}

	for _, buff := range a.Buffs {
		if buff.Stat == DefIgnore && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			defReduction += buff.Value
		}
	}

	totalDef := baseDef*(1+defPct/100-defReduction/100) + flatDef
	totalDef = math.Max(totalDef, 0)

	return 1 - (totalDef / (totalDef + 200.0 + 10.0*float64(c.Level)))
}

func CalcVulnerabilityMultiplier(c Character, e Enemy, a Attack) float64 {
	vuln := 0.0
	for _, buff := range e.Buffs {
		if buff.Stat == Vulnerability && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.AttackTag) {
			vuln += buff.Value
		}
	}
	return 1.0 + vuln/100
}
