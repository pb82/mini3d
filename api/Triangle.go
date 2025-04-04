package api

import (
	"image/color"
)

// Triangle represents a single triangle
type Triangle struct {
	// The list of vertices of the triangle in clockwise order
	Vertices [3]Vector3d

	// The list of texture coordinates
	UVs UVs

	// Optional color
	Color color.Color
}

// RGBA returns the color components of the triangle, implementing the
// `Color` interface
func (t *Triangle) RGBA() (r, g, b, a uint32) {
	return t.Color.RGBA()
}

// Normal computes the normal vector of a triangle
// This is used to determine if a triangle is visible
func (t *Triangle) Normal() Vector3d {
	l1 := t.Vertices[1].Sub(&t.Vertices[0])
	l2 := t.Vertices[2].Sub(&t.Vertices[0])
	normal := l1.Cross(&l2)
	normal.Normalize()
	return normal
}

func (t *Triangle) ScaleW() {
	t.Vertices[0] = t.Vertices[0].Div(t.Vertices[0].W)
	t.Vertices[1] = t.Vertices[1].Div(t.Vertices[1].W)
	t.Vertices[2] = t.Vertices[2].Div(t.Vertices[2].W)
}

// UnpackVertex returns the raw coordinates of a vertex converted so they are usable in the
// drawTriangle function
func (t *Triangle) UnpackVertex(index int) (int, int, float64, float64, float64) {
	return int(t.Vertices[index].X), int(t.Vertices[index].Y), t.UVs[index].U, t.UVs[index].V, t.UVs[index].W
}

// vectorIntersectPlane returns the vertex where the point intersects the plane (screen boundaries)
func vectorIntersectPlane(p, n, start, end *Vector3d, t *float64) Vector3d {
	n.Normalize()
	d := -n.Dot(p)
	ad := start.Dot(n)
	bd := end.Dot(n)
	*t = (-d - ad) / (bd - ad)
	lineStartToEnd := end.Sub(start)
	lineToIntersect := lineStartToEnd.Mul(*t)
	return start.Add(&lineToIntersect)
}

// ClipAgainstPlane splits into two if one or more vertices intersect with screen boundaries
func (t *Triangle) ClipAgainstPlane(p, n *Vector3d, triangleOut1, triangleOut2 *Triangle) int {
	n.Normalize()

	dist := func(point *Vector3d) float64 {
		return n.X*point.X + n.Y*point.Y + n.Z*point.Z - n.Dot(p)
	}

	d0 := dist(&t.Vertices[0])
	d1 := dist(&t.Vertices[1])
	d2 := dist(&t.Vertices[2])

	insidePointCount := 0
	insideUVsCount := 0
	insidePoints := [3]*Vector3d{}
	insideUVs := [3]*VectorUv{}
	outsidePointCount := 0
	outsideUVsCount := 0
	outsidePoints := [3]*Vector3d{}
	outsideUVs := [3]*VectorUv{}

	// Check how many points of the triangle lie inside the
	// screen boundaries
	if d0 >= 0 {
		insidePoints[insidePointCount] = &t.Vertices[0]
		insidePointCount += 1
		insideUVs[insideUVsCount] = &t.UVs[0]
		insideUVsCount += 1
	} else {
		outsidePoints[outsidePointCount] = &t.Vertices[0]
		outsidePointCount += 1
		outsideUVs[outsideUVsCount] = &t.UVs[0]
		outsideUVsCount += 1
	}

	if d1 >= 0 {
		insidePoints[insidePointCount] = &t.Vertices[1]
		insidePointCount += 1
		insideUVs[insideUVsCount] = &t.UVs[1]
		insideUVsCount += 1
	} else {
		outsidePoints[outsidePointCount] = &t.Vertices[1]
		outsidePointCount += 1
		outsideUVs[outsideUVsCount] = &t.UVs[1]
		outsideUVsCount += 1
	}

	if d2 >= 0 {
		insidePoints[insidePointCount] = &t.Vertices[2]
		insidePointCount += 1
		insideUVs[insideUVsCount] = &t.UVs[2]
		insideUVsCount += 1
	} else {
		outsidePoints[outsidePointCount] = &t.Vertices[2]
		outsidePointCount += 1
		outsideUVs[outsideUVsCount] = &t.UVs[2]
		outsideUVsCount += 1
	}

	// No points of the triangle are inside screen boundaries, the
	// triangle is not visible
	if insidePointCount == 0 {
		return 0
	}

	// All points of the triangle are inside the screen boundaries, no
	// clipping is required
	if insidePointCount == 3 {
		*triangleOut1 = *t
		return 1
	}

	// Two points lie outside of screen boundaries. We can clip the triangle into
	// a new, smaller, triangle
	if insidePointCount == 1 && outsidePointCount == 2 {
		triangleOut2.Color = t.Color

		// Keep the inside vertex and UV
		triangleOut1.Vertices[0] = *insidePoints[0]
		triangleOut1.UVs[0] = *insideUVs[0]

		t := 0.0
		triangleOut1.Vertices[1] = vectorIntersectPlane(p, n, insidePoints[0], outsidePoints[0], &t)
		triangleOut1.UVs[1].U = t*(outsideUVs[0].U-insideUVs[0].U) + insideUVs[0].U
		triangleOut1.UVs[1].V = t*(outsideUVs[0].V-insideUVs[0].V) + insideUVs[0].V
		triangleOut1.UVs[1].W = t*(outsideUVs[0].W-insideUVs[0].W) + insideUVs[0].W

		triangleOut1.Vertices[2] = vectorIntersectPlane(p, n, insidePoints[0], outsidePoints[1], &t)
		triangleOut1.UVs[2].U = t*(outsideUVs[1].U-insideUVs[0].U) + insideUVs[0].U
		triangleOut1.UVs[2].V = t*(outsideUVs[1].V-insideUVs[0].V) + insideUVs[0].V
		triangleOut1.UVs[2].W = t*(outsideUVs[1].W-insideUVs[0].W) + insideUVs[0].W

		return 1
	}

	// Two points lie inside of screen boundaries, one outside. Triangle needs to be clipped
	// into two smaller triangles
	if insidePointCount == 2 && outsidePointCount == 1 {
		triangleOut1.Color = t.Color
		triangleOut2.Color = t.Color

		// The first triangle consists of the two inside points and a new
		// point determined by the location where one side of the triangle
		// intersects with the plane
		triangleOut1.Vertices[0] = *insidePoints[0]
		triangleOut1.Vertices[1] = *insidePoints[1]
		triangleOut1.UVs[0] = *insideUVs[0]
		triangleOut1.UVs[1] = *insideUVs[1]

		t := 0.0

		triangleOut1.Vertices[2] = vectorIntersectPlane(p, n, insidePoints[0], outsidePoints[0], &t)
		triangleOut1.UVs[2].U = t*(outsideUVs[0].U-insideUVs[0].U) + insideUVs[0].U
		triangleOut1.UVs[2].V = t*(outsideUVs[0].V-insideUVs[0].V) + insideUVs[0].V
		triangleOut1.UVs[2].W = t*(outsideUVs[0].W-insideUVs[0].W) + insideUVs[0].W

		// The second triangle is composed of one of he inside points, a
		// new point determined by the intersection of the other side of the
		// triangle and the plane, and the newly created point above
		triangleOut2.Vertices[0] = *insidePoints[1]
		triangleOut2.UVs[0] = *insideUVs[1]
		triangleOut2.Vertices[1] = triangleOut1.Vertices[2]
		triangleOut2.UVs[1] = triangleOut1.UVs[2]
		triangleOut2.Vertices[2] = vectorIntersectPlane(p, n, insidePoints[1], outsidePoints[0], &t)
		triangleOut2.UVs[2].U = t*(outsideUVs[0].U-insideUVs[1].U) + insideUVs[1].U
		triangleOut2.UVs[2].V = t*(outsideUVs[0].V-insideUVs[1].V) + insideUVs[1].V
		triangleOut2.UVs[2].W = t*(outsideUVs[0].W-insideUVs[1].W) + insideUVs[1].W

		// Return two newly formed triangles which form a quad
		return 2
	}

	return 0
}

// ColoredTriangleFromMatrix is a utility function to allow the quick creation of a colored
// triangle from a 3x3 matrix using the matrix entries as vertices
func ColoredTriangleFromMatrix(m Matrix3x3, color color.Color) Triangle {
	return Triangle{
		Vertices: [3]Vector3d{
			{X: m[0][0], Y: m[0][1], Z: m[0][2], W: 1},
			{X: m[1][0], Y: m[1][1], Z: m[1][2], W: 1},
			{X: m[2][0], Y: m[2][1], Z: m[2][2], W: 1},
		},
		Color: color,
	}
}

// TexturedTriangleFromMatrix is a utility function to quickly create a textured triangle
func TexturedTriangleFromMatrix(m Matrix3x3, uv Matrix3x2) Triangle {
	return Triangle{
		Vertices: [3]Vector3d{
			{X: m[0][0], Y: m[0][1], Z: m[0][2], W: 1},
			{X: m[1][0], Y: m[1][1], Z: m[1][2], W: 1},
			{X: m[2][0], Y: m[2][1], Z: m[2][2], W: 1},
		},
		UVs: [3]VectorUv{
			{U: uv[0][0], V: uv[0][1]},
			{U: uv[1][0], V: uv[1][1]},
			{U: uv[2][0], V: uv[2][1]},
		},
		Color: nil,
	}
}
