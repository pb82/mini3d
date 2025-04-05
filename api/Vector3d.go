package api

import "math"

// Vector3d represents coordinates in 3d space
type Vector3d struct {
	X, Y, Z float64

	// W represents a scale factor used for perspective correction
	W float64
}

// Copy returns a new `Vector3d` with the same coordinates
func (v *Vector3d) Copy() Vector3d {
	return Vector3d{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
		W: v.W,
	}
}

// Add adds `other` to `v`, producing a new vector
// The `W` component is not affected by this operation
func (v *Vector3d) Add(other *Vector3d) Vector3d {
	return Vector3d{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
		W: v.W,
	}
}

// Sub subtracts `other` from `v`, producing a new vector
// The `W` component is not affected by this operation
func (v *Vector3d) Sub(other *Vector3d) Vector3d {
	return Vector3d{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
		W: v.W,
	}
}

// Mul multiplies the vector with a scalar, producing a new vector
// The `W` component is not affected by this operation
func (v *Vector3d) Mul(scalar float64) Vector3d {
	return Vector3d{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
		W: v.W,
	}
}

// Div divides the vector by a scalar, producing a new vector
// The `W` component is not affected by this operation
func (v *Vector3d) Div(scalar float64) Vector3d {
	return Vector3d{
		X: v.X / scalar,
		Y: v.Y / scalar,
		Z: v.Z / scalar,
		W: v.W,
	}
}

// Dot returns the dot product of `v` and `other`
func (v *Vector3d) Dot(other *Vector3d) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Len returns the length of the vector
func (v *Vector3d) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

// Normalize normalizes the vector by its length (in place)
// The `W` component is not affected by this operation
func (v *Vector3d) Normalize() {
	l := v.Len()
	v.X /= l
	v.Y /= l
	v.Z /= l
}

// Cross computes the cross product of `v` and `other`, producing a new vector
// The `W` component is not affected by this operation
func (v *Vector3d) Cross(other *Vector3d) Vector3d {
	return Vector3d{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
		W: v.W,
	}
}
