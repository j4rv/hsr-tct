package hsrtct

type Character struct {
	ID        uint64
	Name      string
	Level     int
	BaseHp    float64
	BaseAtk   float64
	BaseDef   float64
	BaseSpd   float64
	BaseAggro float64
	Element   Element
	Buffs     []Buff
}

// Will probably be called many times, make it cached in the future
func (c *Character) AllBuffs(lc LightCone, rb RelicBuild) []Buff {
	allBuffs := make([]Buff, len(c.Buffs)+len(lc.Buffs)+len(rb.AsBuffs()))
	copy(allBuffs, c.Buffs)
	copy(allBuffs[len(c.Buffs):], lc.Buffs)
	copy(allBuffs[len(c.Buffs)+len(lc.Buffs):], rb.AsBuffs())
	return allBuffs
}

// TODO cache!
func (c *Character) FinalStatValue(lc LightCone, rb RelicBuild, stat Stat, tag DamageTag, element Element, extraBuffs []Buff) float64 {
	baseValue := 0.0

	switch stat {
	case Hp:
		baseValue = c.BaseHp + lc.BaseHp
	case Atk:
		baseValue = c.BaseAtk + lc.BaseAtk
	case Def:
		baseValue = c.BaseDef + lc.BaseDef
	case Spd:
		baseValue = c.BaseSpd
	case Aggro:
		baseValue = c.BaseAggro
	case CritRate:
		baseValue = 5.0
	case CritDmg:
		baseValue = 50.0
	case EnergyRegenerationRate:
		baseValue = 100.0
	}

	value := baseValue
	allBuffs := c.AllBuffs(lc, rb)
	allBuffs = append(allBuffs, extraBuffs...)
	for _, buff := range allBuffs {
		if !buff.DamageTag.Is(tag) || !buff.Element.Is(element) {
			continue
		}

		switch stat {
		case Hp:
			if buff.Stat == Hp {
				value += buff.Value
			} else if buff.Stat == HpPct {
				value += baseValue * buff.Value / 100
			}
		case Atk:
			if buff.Stat == Atk {
				value += buff.Value
			} else if buff.Stat == AtkPct {
				value += baseValue * buff.Value / 100
			}
		case Def:
			if buff.Stat == Def {
				value += buff.Value
			} else if buff.Stat == DefPct {
				value += baseValue * buff.Value / 100
			}
		case Spd:
			if buff.Stat == Spd {
				value += buff.Value
			} else if buff.Stat == SpdPct {
				value += baseValue * buff.Value / 100
			}
		default:
			if buff.Stat == stat {
				value += buff.Value
			}
		}
	}
	return value
}

type LightCone struct {
	ID      uint64  `json:"id"`
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	BaseHp  float64 `json:"baseHp"`
	BaseAtk float64 `json:"baseAtk"`
	BaseDef float64 `json:"baseDef"`
	Buffs   []Buff  `json:"buffs"`
}
