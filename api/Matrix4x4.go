package api

import "math"

// Matrix4x4 represents a 4 by 4 matrix
type Matrix4x4 [4][4]float64

// Identity4x4 returns a 4 by 4 identity matrix
func Identity4x4() Matrix4x4 {
	return Matrix4x4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Projection4x4 returns a projection matrix to project from 3D into 2D
func Projection4x4(fov, aspectRatio, near, far float64) Matrix4x4 {
	matrix := Matrix4x4{}
	matrix[0][0] = aspectRatio * fov
	matrix[1][1] = fov
	matrix[2][2] = far / (far - near)
	matrix[3][2] = (-far * near) / (far - near)
	matrix[2][3] = 1.0
	matrix[3][3] = 0.0
	return matrix
}

// MulM multiplies the matrix with another matrix of the same order, returning a new matrix
func (m *Matrix4x4) MulM(other *Matrix4x4) Matrix4x4 {
	result := Matrix4x4{}
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			result[r][c] = m[r][0]*other[0][c] + m[r][1]*other[1][c] + m[r][2]*other[2][c] + m[r][3]*other[3][c]
		}
	}
	return result
}

// MulV multiplies the matrix with a vector, returning a new vector
func (m *Matrix4x4) MulV(other *Vector3d) Vector3d {
	result := Vector3d{}
	result.X = other.X*m[0][0] + other.Y*m[1][0] + other.Z*m[2][0] + other.W*m[3][0]
	result.Y = other.X*m[0][1] + other.Y*m[1][1] + other.Z*m[2][1] + other.W*m[3][1]
	result.Z = other.X*m[0][2] + other.Y*m[1][2] + other.Z*m[2][2] + other.W*m[3][2]
	result.W = other.X*m[0][3] + other.Y*m[1][3] + other.Z*m[2][3] + other.W*m[3][3]
	return result
}

// RotateX applies a rotation on the X axis to this matrix (in place)
func (m *Matrix4x4) RotateX(radians float64) {
	m[0][0] = 1
	m[1][1] = math.Cos(radians)
	m[1][2] = math.Sin(radians)
	m[2][1] = -math.Sin(radians)
	m[2][2] = math.Cos(radians)
	m[3][3] = 1
}

// RotateY applies a rotation on the Y axis to this matrix (in place)
func (m *Matrix4x4) RotateY(radians float64) {
	m[0][0] = math.Cos(radians)
	m[0][2] = math.Sin(radians)
	m[2][0] = -math.Sin(radians)
	m[1][1] = 1
	m[2][2] = math.Cos(radians)
	m[3][3] = 1
}

// RotateZ applies a rotation on the Y axis to this matrix (in place)
func (m *Matrix4x4) RotateZ(radians float64) {
	m[0][0] = math.Cos(radians)
	m[0][1] = math.Sin(radians)
	m[1][0] = -math.Sin(radians)
	m[1][1] = math.Cos(radians)
	m[2][2] = 1
	m[3][3] = 1
}

// Translate applies a translation on all three axes to the matrix
func (m *Matrix4x4) Translate(x, y, z float64) {
	m[0][0] = 1
	m[1][1] = 1
	m[2][2] = 1
	m[3][3] = 1
	m[3][0] = x
	m[3][1] = y
	m[3][2] = z
}

// Inverse inverses the matrix, producing a new matrix
func (m *Matrix4x4) Inverse() Matrix4x4 {
	matrix := Matrix4x4{}
	matrix[0][0] = m[0][0]
	matrix[0][1] = m[1][0]
	matrix[0][2] = m[2][0]
	matrix[0][3] = 0.0
	matrix[1][0] = m[0][1]
	matrix[1][1] = m[1][1]
	matrix[1][2] = m[2][1]
	matrix[1][3] = 0.0
	matrix[2][0] = m[0][2]
	matrix[2][1] = m[1][2]
	matrix[2][2] = m[2][2]
	matrix[2][3] = 0.0
	matrix[3][0] = -(m[3][0]*matrix[0][0] + m[3][1]*matrix[1][0] + m[3][2]*matrix[2][0])
	matrix[3][1] = -(m[3][0]*matrix[0][1] + m[3][1]*matrix[1][1] + m[3][2]*matrix[2][1])
	matrix[3][2] = -(m[3][0]*matrix[0][2] + m[3][1]*matrix[1][2] + m[3][2]*matrix[2][2])
	matrix[3][3] = 1.0
	return matrix
}

// PointAt produces a new matrix with camera position and direction applied
func (m *Matrix4x4) PointAt(camera, target, up *Vector3d) {
	newForward := target.Sub(camera)
	newForward.Normalize()
	a := newForward.Mul(up.Dot(&newForward))
	newUp := up.Sub(&a)
	newUp.Normalize()
	newRight := newUp.Cross(&newForward)
	m[0][0] = newRight.X
	m[0][1] = newRight.Y
	m[0][2] = newRight.Z
	m[0][3] = 0.0
	m[1][0] = newUp.X
	m[1][1] = newUp.Y
	m[1][2] = newUp.Z
	m[1][3] = 0.0
	m[2][0] = newForward.X
	m[2][1] = newForward.Y
	m[2][2] = newForward.Z
	m[2][3] = 0.0
	m[3][0] = camera.X
	m[3][1] = camera.Y
	m[3][2] = camera.Z
	m[3][3] = 1.0
}
