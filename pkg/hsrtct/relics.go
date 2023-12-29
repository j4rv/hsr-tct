package hsrtct

type RelicBuild struct {
	ID         uint64
	Relics     [6]Relic
	SubStats   []RelicSubstat
	SetEffects []Buff
}

func (rb *RelicBuild) AsBuffs() []Buff {
	var buffs []Buff
	for _, relic := range rb.Relics {
		buffs = append(buffs, relic.AsBuffs()...)
	}
	for _, substat := range rb.SubStats {
		buffs = append(buffs, substat.AsBuff())
	}
	buffs = append(buffs, rb.SetEffects...)
	return buffs
}

type Relic struct {
	Set      string
	MainStat Stat
	SubStats []RelicSubstat
}

func (r *Relic) AsBuffs() []Buff {
	var buffs []Buff

	mainStatBuff := Buff{
		Stat:  r.MainStat,
		Value: mainStat[r.MainStat],
	}
	buffs = append(buffs, mainStatBuff)

	for _, substat := range r.SubStats {
		subStatBuff := Buff{
			Stat:  substat.Stat,
			Value: substat.Value(),
		}
		buffs = append(buffs, subStatBuff)
	}

	return buffs
}

type RelicSubstat struct {
	Stat     Stat
	Rolls    int
	RollType RollType
	value    float64
}

func (rs *RelicSubstat) AsBuff() Buff {
	return Buff{
		Stat:  rs.Stat,
		Value: rs.Value(),
	}
}

func (rs *RelicSubstat) Value() float64 {
	if rs.value != 0 {
		return rs.value
	}

	value := 0.0
	switch rs.RollType {
	case RollTypeMin:
		value = substatMinRoll[rs.Stat]
	case RollTypeAvg:
		value = substatAvgRoll[rs.Stat]
	case RollTypeMax:
		value = substatMaxRoll[rs.Stat]
	}
	rs.value = value
	return value * float64(rs.Rolls)
}

type RollType string

const (
	RollTypeMin RollType = "Min"
	RollTypeAvg RollType = "Avg"
	RollTypeMax RollType = "Max"
)

// Main stat values
var mainStat map[Stat]float64 = map[Stat]float64{
	Spd:                    25,
	Hp:                     705.6,
	Atk:                    352.8,
	Def:                    352.8,
	HpPct:                  43.20,
	AtkPct:                 43.20,
	DefPct:                 54,
	BreakEffect:            64.80,
	EffectHitRate:          43.20,
	EffectRes:              43.20,
	EnergyRegenerationRate: 19,
	OutgoingHealingBoost:   35,
	DmgBonus:               39,
	CritRate:               32.40,
	CritDmg:                64.80,
}

// Substat values - Min roll
var substatMinRoll map[Stat]float64 = map[Stat]float64{
	Spd:           2.00,
	Hp:            33.87,
	Atk:           16.93,
	Def:           16.93,
	HpPct:         3.46,
	AtkPct:        3.46,
	DefPct:        4.32,
	BreakEffect:   5.18,
	EffectHitRate: 3.46,
	EffectRes:     3.46,
	CritRate:      2.59,
	CritDmg:       5.18,
}

// Substat values - Avg roll
var substatAvgRoll map[Stat]float64 = map[Stat]float64{
	Spd:           2.30,
	Hp:            38.10,
	Atk:           19.05,
	Def:           19.05,
	HpPct:         3.89,
	AtkPct:        3.89,
	DefPct:        4.86,
	BreakEffect:   5.83,
	EffectHitRate: 3.89,
	EffectRes:     3.89,
	CritRate:      2.92,
	CritDmg:       5.83,
}

// Substat values - Max roll
var substatMaxRoll map[Stat]float64 = map[Stat]float64{
	Spd:           2.60,
	Hp:            42.34,
	Atk:           21.17,
	Def:           21.17,
	HpPct:         4.32,
	AtkPct:        4.32,
	DefPct:        5.40,
	BreakEffect:   6.48,
	EffectHitRate: 4.32,
	EffectRes:     4.32,
	CritRate:      3.24,
	CritDmg:       6.48,
}
