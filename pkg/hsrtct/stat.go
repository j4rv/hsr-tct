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

func StatKeys() []Stat {
	return []Stat{
		Hp, Atk, Def, Spd,
		HpPct, AtkPct, DefPct, SpdPct,
		CritRate, CritDmg, OutgoingHealingBoost, EffectHitRate, EffectRes,
		EnergyRegenerationRate, BreakEffect, DefIgnore, DefShred, Aggro,
		DmgBonus, ElementalRes, ResShred, ResPen, Vulnerability,
	}
}

type Buff struct {
	Stat      Stat      `json:"stat"`
	Value     float64   `json:"value"`
	DamageTag DamageTag `json:"damageTag"`
	Element   Element   `json:"element"`
}

func (b Buff) String() string {
	prettyStat := strings.Replace(string(b.Stat), "Pct", "%", 1)
	valueSuffix := ""
	if !(b.Stat == Hp || b.Stat == Atk || b.Stat == Def || b.Stat == Spd || b.Stat == Aggro) {
		valueSuffix = "%"
	}
	result := fmt.Sprintf("%.1f%s %s", b.Value, valueSuffix, prettyStat)

	if b.DamageTag != "" {
		result += fmt.Sprintf("(%s)", b.DamageTag)
	}

	if b.Element != "" {
		result += fmt.Sprintf("(%s)", b.Element)
	}

	return result
}
