package function

// Space : The abstract world
type Space struct {
	plot        func(v Vector)
	dimension   int
	raw         []Vector
	basicVector []Vector
	subspace    []SubSpace
	pointer     Vector
	density     float64
}



type Vector []float64

type SubSpace struct {
	begin Vector
	end   Vector
}

// CLear Remove all dots from space
func CLear(space *Space) {
	space.raw = nil
}

func (s *Space) Plot(v Vector) {
	s.plot(v)
}

func (s *Space) GetPointer() Vector {
	return s.pointer
}

func (s *Space) SetPointer(v Vector) {
	s.pointer = v
}

func (s *Space) PlotArray(v []Vector) {
	for _, vv := range v {
		s.Plot(vv)
	}
}

func NewSpace() *Space {
	return & Space {
		basicVector: []Vector{
			{0, 0, 1}, {0, 1, 0}, {1, 0, 0},
		},
		pointer: []float64{0, 0, 0},
	}
}