package api

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type Mesh []Triangle

type Mesh struct {
	Triangles []Triangle
}

// Copy returns a new mesh with copies of the same triangles
func (m *Mesh) Copy() *Mesh {
	duplicate := &Mesh{}
	for _, t := range m.Triangles {
		duplicate.Triangles = append(duplicate.Triangles, t.Copy())
	}
	return duplicate
}

// SetMeshPositionRelative move the whole mesh to a new position given relative coordinates
func (m *Mesh) SetMeshPositionRelative(dx, dy, dz float64) {
	for i := range m.Triangles {
		m.Triangles[i].SetTrianglePositionRelative(dx, dy, dz)
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

	result := &Mesh{}

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
			if currentLine[1] == 't' {
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
				result.Triangles = append(result.Triangles, triangles...)
			} else {
				triangles, err := parseFaceWithTexture(currentLine[2:], lineNumber)
				if err != nil {
					return nil, err
				}
				result.Triangles = append(result.Triangles, triangles...)
			}
		}
	}

	return result, nil
}
