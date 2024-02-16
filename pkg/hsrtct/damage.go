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

func AllDamageTags() []DamageTag {
	return []DamageTag{Basic, Skill, Ultimate, FollowUp, Dot}
}

type AttackAOE string

const (
	Single            AttackAOE = ""
	Blast             AttackAOE = "Blast"
	All               AttackAOE = "All"
	EvenlyDistributed AttackAOE = "EvenlyDistributed"
)

func ParseAttackAOE(s string) (AttackAOE, error) {
	switch s {
	case "":
	case "Single":
		return Single, nil
	case "Blast":
		return Blast, nil
	case "All":
		return All, nil
	case "EvenlyDistributed":
		return EvenlyDistributed, nil
	}
	return Single, fmt.Errorf("invalid attack AOE: %s", s)
}

// Attack is an attack that can be used by a character.
// If AOE is Single: will use Multiplier and MultiplierSplash will be ignored.
// If AOE is Blast, will use Multiplier for the focused enemy and MultiplierSplash for its neighbors.
// If AOE is All, will use Multiplier for all enemies.
// If AOE is EvenlyDistributed, will use a percentage of Multiplier on each enemy, evenly distributed.
type Attack struct {
	ID               uint64
	Name             string
	ScalingStat      Stat
	Multiplier       float64
	MultiplierSplash float64
	Element          Element
	DamageTag        DamageTag
	AttackAOE        AttackAOE
	Buffs            []Buff
}

type Scenario struct {
	ID           uint64
	Name         string
	Notes        string
	Character    Character
	LightCone    LightCone
	RelicBuild   RelicBuild
	Enemies      []Enemy
	FocusedEnemy int
	Attacks      map[*Attack]float64
}

type ScenarioResult struct {
	TotalDmg     float64
	Explanations []string
}

func CalcAvgDmgScenario(s Scenario) (ScenarioResult, error) {
	totalDmg := 0.0
	explanations := []string{}
	for attack, mult := range s.Attacks {
		switch attack.AttackAOE {

		case Single:
			dmg, exp, err := CalcAvgDamage(s.Character, s.LightCone, s.RelicBuild, s.Enemies[s.FocusedEnemy], *attack, false)
			if err != nil {
				return ScenarioResult{}, err
			}
			totalDmg += dmg * mult
			explanations = append(explanations, fmt.Sprintf("%s on %s:\nDamage: %f\n\n%s", attack.Name, s.Enemies[s.FocusedEnemy].Name, dmg, exp))

		case Blast:
			avgDmg, exp, err := CalcAvgDamage(s.Character, s.LightCone, s.RelicBuild, s.Enemies[s.FocusedEnemy], *attack, false)
			if err != nil {
				return ScenarioResult{}, err
			}
			totalDmg += avgDmg * mult
			explanations = append(explanations, fmt.Sprintf("%s on %s (center):\nDamage: %f\n\n%s", attack.Name, s.Enemies[s.FocusedEnemy].Name, avgDmg, exp))
			if s.FocusedEnemy-1 >= 0 {
				splashDmg, exp, err := CalcAvgDamage(s.Character, s.LightCone, s.RelicBuild, s.Enemies[s.FocusedEnemy-1], *attack, true)
				if err != nil {
					return ScenarioResult{}, err
				}
				explanations = append(explanations, fmt.Sprintf("%s on %s (left):\nDamage: %f\n\n%s", attack.Name, s.Enemies[s.FocusedEnemy-1].Name, splashDmg, exp))
				totalDmg += splashDmg * mult
			}
			if s.FocusedEnemy+1 < len(s.Enemies) {
				splashDmg, exp, err := CalcAvgDamage(s.Character, s.LightCone, s.RelicBuild, s.Enemies[s.FocusedEnemy+1], *attack, true)
				if err != nil {
					return ScenarioResult{}, err
				}
				explanations = append(explanations, fmt.Sprintf("%s on %s (right):\nDamage: %f\n\n%s", attack.Name, s.Enemies[s.FocusedEnemy+1].Name, splashDmg, exp))
				totalDmg += splashDmg * mult
			}

		case All, EvenlyDistributed:
			for _, enemy := range s.Enemies {
				dmg, exp, err := CalcAvgDamage(s.Character, s.LightCone, s.RelicBuild, enemy, *attack, false)
				if err != nil {
					return ScenarioResult{}, err
				}
				if attack.AttackAOE == EvenlyDistributed {
					dmg /= float64(len(s.Enemies))
				}
				explanations = append(explanations, fmt.Sprintf("%s on %s:\nDamage: %f\n\n%s", attack.Name, enemy.Name, dmg, exp))
				totalDmg += dmg * mult
			}
		}
	}
	return ScenarioResult{totalDmg, explanations}, nil
}

func CalcAvgDamage(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack, isSplash bool) (float64, string, error) {
	baseDamage, err := CalcBaseDamage(c, lc, rb, e, a, isSplash)
	if err != nil {
		return 0, "", err
	}
	critMult := CalcAvgCritMultiplier(c, lc, rb, e, a)
	dmgBonusMult := CalcDmgBonusMult(c, lc, rb, e, a)
	resMult := CalcResistanceMultiplier(c, lc, rb, e, a)
	defMult := CalcDefenseMultiplier(c, lc, rb, e, a)
	vulnMult := CalcVulnerabilityMultiplier(e, a)
	dmgReductionMult := CalcDmgReductionMultiplier(e, a)

	stats, err := CharacterStats(c, lc, rb, e, a)
	if err != nil {
		return 0, "", err
	}

	explanation := fmt.Sprintf(
		"Base Damage: %.2f\n"+
			"Crit Multiplier: %.2f\n"+
			"Damage Bonus Multiplier: %.2f\n"+
			"Resistance Multiplier: %.2f\n"+
			"Defense Multiplier: %.2f\n"+
			"Vulnerability Multiplier: %.2f\n"+
			"Damage Reduction Multiplier: %.2f\n\n"+
			"Stats:",
		baseDamage, critMult, dmgBonusMult, resMult, defMult, vulnMult, dmgReductionMult)
	for _, stat := range AllStats() {
		value, ok := stats[stat]
		if !ok {
			continue
		}
		explanation += fmt.Sprintf("\n%s: %.2f", stat, value)
	}

	return baseDamage * critMult * dmgBonusMult * resMult * defMult * vulnMult * dmgReductionMult, explanation, nil
}

func CharacterStats(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack) (map[Stat]float64, error) {
	stats := make(map[Stat]float64)

	for _, stat := range AllStats() {
		value := c.FinalStatValue(lc, rb, stat, a.DamageTag, a.Element, nil)
		if value != 0 {
			stats[stat] = value
		}
	}

	return stats, nil
}

func CalcBaseDamage(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack, isSplash bool) (float64, error) {
	baseDamage := 0.0
	mult := a.Multiplier
	if isSplash {
		mult = a.MultiplierSplash
	}
	switch a.ScalingStat {
	case Hp:
		baseDamage = c.FinalStatValue(lc, rb, Hp, a.DamageTag, a.Element, a.Buffs) * mult / 100
	case Atk:
		baseDamage = c.FinalStatValue(lc, rb, Atk, a.DamageTag, a.Element, a.Buffs) * mult / 100
	case Def:
		baseDamage = c.FinalStatValue(lc, rb, Def, a.DamageTag, a.Element, a.Buffs) * mult / 100
	default:
		return 0, ErrInvalidScalingStat
	}
	return baseDamage, nil
}

func CalcAvgCritMultiplier(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack) float64 {
	if a.DamageTag == Dot {
		return 1
	}
	critRate := c.FinalStatValue(lc, rb, CritRate, a.DamageTag, a.Element, a.Buffs)
	critDamage := c.FinalStatValue(lc, rb, CritDmg, a.DamageTag, a.Element, a.Buffs)
	critRate = math.Min(critRate, 100)
	return 1 + (critRate / 100 * critDamage / 100)
}

func CalcDmgBonusMult(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack) float64 {
	return 1 + c.FinalStatValue(lc, rb, DmgBonus, a.DamageTag, a.Element, a.Buffs)/100
}

func CalcResistanceMultiplier(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack) float64 {
	res := 0.0

	for _, buff := range e.Buffs {
		if buff.Stat == ElementalRes && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			res += buff.Value
		}
		if buff.Stat == ResShred && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			res -= buff.Value
		}
	}

	for _, buff := range c.AllBuffs(lc, rb) {
		if buff.Stat == ResPen && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			res -= buff.Value
		}
	}

	for _, buff := range a.Buffs {
		if buff.Stat == ResPen && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			res -= buff.Value
		}
	}

	return 1.0 - res/100
}

func CalcDefenseMultiplier(c Character, lc LightCone, rb RelicBuild, e Enemy, a Attack) float64 {
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
		if buff.Stat == DefShred && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			defReduction += buff.Value
		}
	}

	for _, buff := range c.AllBuffs(lc, rb) {
		if buff.Stat == DefIgnore && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			defReduction += buff.Value
		}
	}

	for _, buff := range a.Buffs {
		if buff.Stat == DefIgnore && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			defReduction += buff.Value
		}
	}

	totalDef := baseDef*(1+defPct/100-defReduction/100) + flatDef
	totalDef = math.Max(totalDef, 0)

	return 1 - (totalDef / (totalDef + 200.0 + 10.0*float64(c.Level)))
}

func CalcVulnerabilityMultiplier(e Enemy, a Attack) float64 {
	vuln := 0.0
	for _, buff := range e.Buffs {
		if buff.Stat == Vulnerability && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			vuln += buff.Value
		}
	}
	return 1.0 + vuln/100
}

func CalcDmgReductionMultiplier(e Enemy, a Attack) float64 {
	dmgReduction := 0.0
	for _, buff := range e.Buffs {
		if buff.Stat == DmgReduction && buff.Element.Is(a.Element) && buff.DamageTag.Is(a.DamageTag) {
			dmgReduction += buff.Value
		}
	}
	return 1.0 - dmgReduction/100
}
