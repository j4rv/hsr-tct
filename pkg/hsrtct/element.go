package hsrtct

type Element string

const (
	AnyElement Element = ""
	Ice        Element = "Ice"
	Wind       Element = "Wind"
	Fire       Element = "Fire"
	Imaginary  Element = "Imaginary"
	Lightning  Element = "Lightning"
	Quantum    Element = "Quantum"
	Physical   Element = "Physical"
)

func (e Element) Is(element Element) bool {
	return e == element || e == AnyElement || element == AnyElement
}
