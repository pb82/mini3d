package api

// VectorUv represents a single texture coordinate
type VectorUv struct {
	U, V float64

	// W represents a scale factor
	W float64
}

// UVs represents a set of texture coordinates
type UVs [3]VectorUv

// Copy copies a set of UV vectors
func (u *UVs) Copy() UVs {
	return UVs{
		u[0],
		u[1],
		u[2],
	}
}

// ScaleW scales the U and W coordinates by the `W` component
func (u *UVs) ScaleW(t *Triangle) {
	for n := 0; n < len(u); n++ {
		u[n].U = u[n].U / t.Vertices[n].W
		u[n].V = u[n].V / t.Vertices[n].W
		u[n].W = 1 / t.Vertices[n].W
	}
}
