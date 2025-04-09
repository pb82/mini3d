package api

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Mesh struct {
	triangles []Triangle

	// Rotation around origin
	rotX Matrix4x4
	rotY Matrix4x4
	rotZ Matrix4x4

	// Mesh translation
	trans Matrix4x4

	// Rotation around arbitrary points
	rotXAround Matrix4x4
	rotYAround Matrix4x4
	rotZAround Matrix4x4

	// World matrix to apply all transformations to
	world Matrix4x4

	// Internal helper matrices for multistep transformations, e.g.
	// translate, rotate, translate
	helperMatrix1 Matrix4x4
	helperMatrix2 Matrix4x4
	helperMatrix3 Matrix4x4

	// Keep track of the position and the dimensions of the bounding box
	minX, minY, minZ float64
	maxX, maxY, maxZ float64
}

// NewMesh creates a new mesh instance, with sane initial values
// The mesh will be positioned at origin (0,0,0)
func NewMesh() *Mesh {
	mesh := &Mesh{}
	mesh.rotX = Identity4x4()
	mesh.rotY = Identity4x4()
	mesh.rotZ = Identity4x4()
	mesh.rotXAround = Identity4x4()
	mesh.rotYAround = Identity4x4()
	mesh.rotZAround = Identity4x4()
	mesh.world = Identity4x4()
	mesh.Translate(0, 0, 0)
	return mesh
}

// updateBoundingBox establishes a bounding box of minimum and maximum coordinates
// in every direction
func (m *Mesh) updateBoundingBox(v *Vector3d) {
	if v.X < m.minX {
		m.minX = v.X
	}
	if v.Y < m.minY {
		m.minY = v.Y
	}
	if v.Z < m.minZ {
		m.minZ = v.Z
	}
	if v.X > m.maxX {
		m.maxX = v.X
	}
	if v.Y > m.maxY {
		m.maxY = v.Y
	}
	if v.Z > m.maxZ {
		m.maxZ = v.Z
	}
}

// AddTriangle adds a single triangle to the mesh
func (m *Mesh) AddTriangle(triangle Triangle) {
	m.triangles = append(m.triangles, triangle)
	for _, v := range triangle.Vertices {
		m.updateBoundingBox(&v)
	}
}

// AddTriangles adds a list of triangles to the mesh
func (m *Mesh) AddTriangles(triangles []Triangle) {
	for _, triangle := range triangles {
		m.AddTriangle(triangle)
	}
}

// GetOrigin return the origin coordinate of the bounding box of the mesh
func (m *Mesh) GetOrigin() Vector3d {
	return Vector3d{
		X: m.minX,
		Y: m.minY,
		Z: m.minZ,
	}
}

// GetCenter returns the center of the bounding box
func (m *Mesh) GetCenter() Vector3d {
	b := m.GetBoundingBox()
	o := m.GetOrigin()
	return Vector3d{
		X: -o.X - b.X/2,
		Y: -o.Y - b.Y/2,
		Z: -o.Z - b.Z/2,
	}
}

// GetBoundingBox return the dimensions of the bounding nox of the mesh
func (m *Mesh) GetBoundingBox() Vector3d {
	v := Vector3d{
		X: m.maxX - m.minX,
		Y: m.maxY - m.minY,
		Z: m.maxZ - m.minZ,
	}
	return v
}

func (m *Mesh) reset() {
	m.minX = 99999
	m.maxX = -99999
	m.minY = 99999
	m.maxX = -99999
	m.minZ = 99999
	m.maxZ = -99999
}

func (m *Mesh) MoveRelative(dx, dy, dz float64) {
	m.reset()
	for i := range m.triangles {
		for v := range m.triangles[i].Vertices {
			m.triangles[i].Vertices[v].X += dx
			m.triangles[i].Vertices[v].Y += dy
			m.triangles[i].Vertices[v].Z += dz
			m.updateBoundingBox(&m.triangles[i].Vertices[v])
		}
	}
}

// Translate translates the mesh to an absolute position
func (m *Mesh) Translate(x, y, z float64) {
	m.trans.Translate(x, y, z)
	m.minX += x
	m.maxX += x
	m.minY += y
	m.maxY += y
	m.minZ += z
	m.maxZ += z
}

// RotateXAroundOrigin rotates all meshes around origin on the X axis
func (m *Mesh) RotateXAroundOrigin(radians float64) {
	m.rotX.RotateX(radians)
}

// RotateYAroundOrigin rotates the mesh around origin on the Y axis
func (m *Mesh) RotateYAroundOrigin(radians float64) {
	m.rotY.RotateY(radians)
}

// RotateZAroundOrigin rotates the mesh around origin on the Z axis
func (m *Mesh) RotateZAroundOrigin(radians float64) {
	m.rotZ.RotateZ(radians)
}

// RotateXAround rotates the mesh around a given point on the X axis
func (m *Mesh) RotateXAround(radians float64, v *Vector3d) {
	m.helperMatrix1 = Identity4x4()
	m.helperMatrix1.Translate(0, v.Y, v.Z)
	m.helperMatrix2 = Identity4x4()
	m.helperMatrix2.RotateX(radians)
	m.helperMatrix3 = Identity4x4()
	m.helperMatrix3.Translate(0, -v.Y, -v.Z)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix2)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix3)
	m.rotXAround = m.helperMatrix1
}

// RotateYAround rotates the mesh around a given point on the X axis
func (m *Mesh) RotateYAround(radians float64, v *Vector3d) {
	m.helperMatrix1 = Identity4x4()
	m.helperMatrix1.Translate(v.X, 0, v.Z)
	m.helperMatrix2 = Identity4x4()
	m.helperMatrix2.RotateY(radians)
	m.helperMatrix3 = Identity4x4()
	m.helperMatrix3.Translate(-v.X, 0, -v.Z)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix2)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix3)
	m.rotYAround = m.helperMatrix1
}

// RotateZAround rotates the mesh around a given point on the X axis
func (m *Mesh) RotateZAround(radians float64, v *Vector3d) {
	m.helperMatrix1 = Identity4x4()
	m.helperMatrix1.Translate(v.X, v.Y, 0)
	m.helperMatrix2 = Identity4x4()
	m.helperMatrix2.RotateZ(radians)
	m.helperMatrix3 = Identity4x4()
	m.helperMatrix3.Translate(-v.X, -v.Y, 0)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix2)
	m.helperMatrix1 = m.helperMatrix1.MulM(&m.helperMatrix3)
	m.rotZAround = m.helperMatrix1
}

// Copy returns a new mesh with copies of the same triangles
func (m *Mesh) Copy() *Mesh {
	duplicate := NewMesh()
	for _, t := range m.triangles {
		duplicate.triangles = append(duplicate.triangles, t.Copy())
	}
	return duplicate
}

// SetMeshPositionRelative move the whole mesh to a new position given relative coordinates
func (m *Mesh) SetMeshPositionRelative(dx, dy, dz float64) {
	for i := range m.triangles {
		m.triangles[i].SetTrianglePositionRelative(dx, dy, dz)
	}
}

// LoadWavefrontObj implements rudimentary Wavefront obj file format support
func LoadWavefrontObj(filename string) (*Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	result := NewMesh()

	// Arrays to contain all vertices and uvs found in the file
	// They will later be referenced from face information
	var vertices []Vector3d
	var uvs []VectorUv

	parseUV := func(line string, lineNumber int) (*VectorUv, error) {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid vt line: '%s' in line %d", line, lineNumber)
		}

		u, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float in vt line: '%s' in line %d", line, lineNumber)
		}

		v, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float in vt line: '%s' in line %d", line, lineNumber)
		}

		return &VectorUv{u, v, 1}, nil
	}

	parseVertex := func(line string, lineNumber int) (*Vector3d, error) {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid v line: '%s' in line %d", line, lineNumber)
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
		}

		z, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
		}

		return &Vector3d{x, y, z, 1}, nil
	}

	parseFaceNoTexture := func(line string, lineNumber int) ([]Triangle, error) {
		triangles := []Triangle{}
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			fa, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			fb, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			fc, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[fa-1], vertices[fb-1], vertices[fc-1]},
			})
		} else if len(parts) == 4 {
			fa, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			fb, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			fc, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			fd, err := strconv.ParseInt(parts[3], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float in v line: '%s' in line %d", line, lineNumber)
			}

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[fa-1], vertices[fb-1], vertices[fc-1]},
			})

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[fa-1], vertices[fc-1], vertices[fd-1]},
			})
		} else {
			return nil, fmt.Errorf("invalid face line: '%s' in line %d", line, lineNumber)
		}

		return triangles, nil
	}

	parseVertexWithTexture := func(part string, lineNumber int) (int64, int64, error) {
		parts := strings.Split(part, "/")
		if len(parts) == 3 || len(parts) == 2 {
			vertexIndex, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				return 0, 0, err
			}
			uvIndex, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return 0, 0, err
			}
			return vertexIndex, uvIndex, nil
		} else {
			return 0, 0, fmt.Errorf("invalid face line: '%s' in line %d", part, lineNumber)
		}
	}

	parseFaceWithTexture := func(line string, lineNumber int) ([]Triangle, error) {
		var triangles []Triangle
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			vertexA, uvA, err := parseVertexWithTexture(parts[0], lineNumber)
			if err != nil {
				return nil, err
			}

			vertexB, uvB, err := parseVertexWithTexture(parts[1], lineNumber)
			if err != nil {
				return nil, err
			}

			vertexC, uvC, err := parseVertexWithTexture(parts[2], lineNumber)
			if err != nil {
				return nil, err
			}

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[vertexA-1], vertices[vertexB-1], vertices[vertexC-1]},
				UVs:      [3]VectorUv{uvs[uvA-1], uvs[uvB-1], uvs[uvC-1]},
			})
		} else if len(parts) == 4 {
			fa, uva, err := parseVertexWithTexture(parts[0], lineNumber)
			if err != nil {
				return nil, err
			}

			fb, uvb, err := parseVertexWithTexture(parts[1], lineNumber)
			if err != nil {
				return nil, err
			}

			fc, uvc, err := parseVertexWithTexture(parts[2], lineNumber)
			if err != nil {
				return nil, err
			}

			fd, uvd, err := parseVertexWithTexture(parts[3], lineNumber)
			if err != nil {
				return nil, err
			}

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[fa-1], vertices[fb-1], vertices[fc-1]},
				UVs:      [3]VectorUv{uvs[uva-1], uvs[uvb-1], uvs[uvc-1]},
			})

			triangles = append(triangles, Triangle{
				Vertices: [3]Vector3d{vertices[fa-1], vertices[fc-1], vertices[fd-1]},
				UVs:      [3]VectorUv{uvs[uva-1], uvs[uvc-1], uvs[uvd-1]},
			})
		} else {
			return nil, fmt.Errorf("invalid face line: '%s' in line %d", line, lineNumber)
		}

		return triangles, nil
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		currentLine := scanner.Text()

		if strings.TrimSpace(currentLine) == "" {
			continue
		}

		if currentLine[0] == 'v' {
			if currentLine[1] == 'n' {
				// `vn` vertex normal. Not supported.
				continue
			} else if currentLine[1] == 't' {
				// `vt` (vertex texture)
				uv, err := parseUV(currentLine[3:], lineNumber)
				if err != nil {
					return nil, err
				}
				uvs = append(uvs, *uv)
			} else {
				// `v` (vertex)
				vertex, err := parseVertex(currentLine[2:], lineNumber)
				if err != nil {
					return nil, err
				}
				vertices = append(vertices, *vertex)
			}
		} else if currentLine[0] == 'f' {
			// The texture index can be appended after the vertex index, separated by a slash
			hasTextures := strings.Contains(currentLine, "/")
			if !hasTextures {
				triangles, err := parseFaceNoTexture(currentLine[2:], lineNumber)
				if err != nil {
					return nil, err
				}
				result.AddTriangles(triangles)
			} else {
				triangles, err := parseFaceWithTexture(currentLine[2:], lineNumber)
				if err != nil {
					return nil, err
				}
				result.AddTriangles(triangles)
			}
		}
	}

	return result, nil
}
