package api

// StandardCube returns a standard cube with the dimensions 1x1x1 at origin 0/0/0
// UV coordinates are set to map the whole texture from 0/0 to 1/1
func StandardCube() *Mesh {
	mesh := &Mesh{}

	// Side 1
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 0, 0},
			{0, 1, 0},
			{1, 1, 0},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 0, 0},
			{1, 1, 0},
			{1, 0, 0},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 2
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 0},
			{1, 1, 0},
			{1, 1, 1},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 1},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 3
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 1},
			{1, 1, 1},
			{0, 1, 1},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 1},
			{0, 1, 1},
			{0, 0, 1},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 4
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 0, 1},
			{0, 1, 1},
			{0, 1, 0},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 0, 1},
			{0, 1, 0},
			{0, 0, 0},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 5
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 1, 0},
			{0, 1, 1},
			{1, 1, 1},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{0, 1, 0},
			{1, 1, 1},
			{1, 1, 0},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 6
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 1},
			{0, 0, 1},
			{0, 0, 0},
		},
		Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh.Triangles = append(mesh.Triangles, TexturedTriangleFromMatrix(
		Matrix3x3{
			{1, 0, 1},
			{0, 0, 0},
			{1, 0, 0},
		},
		Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	return mesh
}
