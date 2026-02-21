package exene

type Pair[A any, B any] struct {
	p1 A
	p2 B
}

func NewPair[A any, B any](v1 A, v2 B) Pair[A, B] {
	return Pair[A, B]{v1, v2}
}

func (p Pair[A, B]) First() A {
	return p.p1
}

func (p Pair[A, B]) Second() B {
	return p.p2
}

func (p Pair[A, B]) Get() (A, B) {
	return p.p1, p.p2
}
