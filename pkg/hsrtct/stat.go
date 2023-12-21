package hsrtct

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
	ID        uint
	Stat      Stat
	Value     float64
	DamageTag AttackTag
	Element   Element
}
