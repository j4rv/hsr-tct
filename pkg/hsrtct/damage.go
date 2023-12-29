package hsrtct

import (
	"errors"
	"fmt"
	"math"
)

var ErrInvalidScalingStat = errors.New("invalid scaling stat")

type AttackTag string

const (
	AnyAttack AttackTag = ""
	Basic     AttackTag = "Basic"
	Skill     AttackTag = "Skill"
	Ultimate  AttackTag = "Ultimate"
	FollowUp  AttackTag = "FollowUp"
	Dot       AttackTag = "Dot"
)

func (d AttackTag) Is(tag AttackTag) bool {
	return d == tag || d == AnyAttack || tag == AnyAttack
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
	ID               string
	Name             string
	ScalingStat      Stat
	Multiplier       float64
	MultiplierSplash float64
	Element          Element
	AttackTag        AttackTag
	AttackAOE        AttackAOE
	Buffs            []Buff
}

type Scenario struct {
	ID           string
	Name         string
	Character    string
	Enemies      []string
	FocusedEnemy string
	Attacks      map[string]float64
}

func CalcAvgDamageST(c Character, e Enemy, a Attack) (float64, error) {
	if a.AttackAOE != Single {
		return 0, errors.New("expected aoe: single")
	}
	baseDamage, err := CalcBaseDamage(c, e, a)
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

func ExplainDamageST(c Character, e Enemy, a Attack) (string, error) {
	baseDamage, err := CalcBaseDamage(c, e, a)
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

func CalcBaseDamage(c Character, e Enemy, a Attack) (float64, error) {
	baseDamage := 0.0
	switch a.ScalingStat {
	case Hp:
		baseDamage = c.FinalStatValue(Hp, a.AttackTag, a.Element, a.Buffs) * a.Multiplier / 100
	case Atk:
		baseDamage = c.FinalStatValue(Atk, a.AttackTag, a.Element, a.Buffs) * a.Multiplier / 100
	case Def:
		baseDamage = c.FinalStatValue(Def, a.AttackTag, a.Element, a.Buffs) * a.Multiplier / 100
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
	}

	for _, debuff := range e.Debuffs {
		if debuff.Stat == ResShred && debuff.Element.Is(a.Element) && debuff.DamageTag.Is(a.AttackTag) {
			res -= debuff.Value
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
	}

	for _, debuff := range e.Debuffs {
		if debuff.Stat == DefShred && debuff.Element.Is(a.Element) && debuff.DamageTag.Is(a.AttackTag) {
			defReduction += debuff.Value
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
	for _, debuff := range e.Debuffs {
		if debuff.Stat == Vulnerability && debuff.Element.Is(a.Element) && debuff.DamageTag.Is(a.AttackTag) {
			vuln += debuff.Value
		}
	}
	return 1.0 + vuln/100
}
