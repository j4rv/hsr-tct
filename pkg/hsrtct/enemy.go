package hsrtct

type Enemy struct {
	ID      uint64
	Name    string
	Level   int
	Buffs   []Buff
	Debuffs []Buff
}
