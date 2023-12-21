package hsrtct

import (
	"fmt"
	"strings"
)

type Stat string

const (
	Hp                     Stat = "Hp"
	Atk                    Stat = "Atk"
	Def                    Stat = "Def"
	Spd                    Stat = "Spd"
	HpPct                  Stat = "HpPct"
	AtkPct                 Stat = "AtkPct"
	DefPct                 Stat = "DefPct"
	SpdPct                 Stat = "SpdPct"
	CritRate               Stat = "CritRate"
	CritDmg                Stat = "CritDmg"
	OutgoingHealingBoost   Stat = "OutgoingHealingBoost"
	EffectHitRate          Stat = "EffectHitRate"
	EffectRes              Stat = "EffectRes"
	EnergyRegenerationRate Stat = "EnergyRegenerationRate"
	BreakEffect            Stat = "BreakEffect"
	DefIgnore              Stat = "DefIgnore"
	DefShred               Stat = "DefShred"
	Aggro                  Stat = "Aggro"
	DmgBonus               Stat = "DmgBonus"
	ElementalRes           Stat = "ElementalRes"
	ResShred               Stat = "ResShred"
	ResPen                 Stat = "ResPen"
	Vulnerability          Stat = "Vulnerability"
)

type Buff struct {
	Stat      Stat
	Value     float64
	DamageTag AttackTag
	Element   Element
}

func (b Buff) String() string {
	prettyStat := strings.Replace(string(b.Stat), "Pct", "%", 1)
	valueSuffix := ""
	if !(b.Stat == Hp || b.Stat == Atk || b.Stat == Def || b.Stat == Spd || b.Stat == Aggro) {
		valueSuffix = "%"
	}
	result := fmt.Sprintf("%.1f%s %s", b.Value, valueSuffix, prettyStat)

	if b.DamageTag != "" {
		result += fmt.Sprintf(", DamageTag: %s", b.DamageTag)
	}

	if b.Element != "" {
		result += fmt.Sprintf(", Element: %s", b.Element)
	}

	return result
}
