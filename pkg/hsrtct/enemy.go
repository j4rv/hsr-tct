package hsrtct

type Enemy struct {
	ID      string
	Name    string
	Level   int
	Buffs   []Buff
	Debuffs []Buff
}
