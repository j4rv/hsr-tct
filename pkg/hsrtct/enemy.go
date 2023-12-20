package hsrtct

type Enemy struct {
	ID      uint
	Name    string
	Level   int
	Buffs   []Buff
	Debuffs []Buff
}
