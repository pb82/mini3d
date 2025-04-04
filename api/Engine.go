package api

import (
	"image/color"
	"math"
	"time"
)

type UserData interface{}
type DrawHook func(x, y int, c color.Color, userData UserData)
type Engine struct {
	// Internal viewport dimensions
	w, h int
	W, H float64

	// List of meshes to render
	meshes []Mesh

	// Rotation of all meshes on the X axis
	rotX Matrix4x4

	// Rotation of all meshes on the X axis
	rotY Matrix4x4

	// Rotation of all meshes on the X axis
	rotZ Matrix4x4

	// Translation of all meshes on all axes
	trans Matrix4x4

	// World matrix to apply all transformations to
	world Matrix4x4

	// Camera yaw angle in radians (left/right)
	yaw float64

	// Camera pitch angle in radians (up/down)
	pitch float64

	// Camera position
	camera Vector3d

	// View matrix
	view Matrix4x4

	// Vector pointing to the current forward direction of the camera
	direction Vector3d

	// Projection matrix to project from 3D into 2D
	projection Matrix4x4

	// Optional texture atlas
	// If this is not set, triangles must have a defined color
	textureAtlas TextureAtlas

	// depthBuffer helps to avoid drawing pixels that have already been filled
	// this implementation does not allow for transparent materials
	depthBuffer *DepthBuffer

	// Hooks - callback functions to be defined by the user of the library
	drawPixel DrawHook

	// yOrigin set the position of the (0/0) coordinate
	yOrigin YOrigin

	// Metrics contains performance indicators
	Metrics Metrics
}

func timestamp() float64 {
	return float64(time.Now().UnixMilli())
}

// AddMesh adds a mesh to the engine in order to be rendered
func (e *Engine) AddMesh(mesh Mesh) {
	e.meshes = append(e.meshes, mesh)
}

// Translate translates all meshes
func (e *Engine) Translate(x, y, z float64) {
	e.trans.Translate(x, y, z)
}

// RotateX rotates all meshes on the X axis
func (e *Engine) RotateX(radians float64) {
	e.rotX.RotateX(radians)
}

// RotateY rotates all meshes on the Y axis
func (e *Engine) RotateY(radians float64) {
	e.rotY.RotateY(radians)
}

// RotateZ rotates all meshes on the Z axis
func (e *Engine) RotateZ(radians float64) {
	e.rotZ.RotateZ(radians)
}

// SetCameraPosition sets the camera to the given position
func (e *Engine) SetCameraPosition(x, y, z float64) {
	e.camera.X = x
	e.camera.Y = y
	e.camera.Z = z
}

// MoveCamera the camera by the given offsets
func (e *Engine) MoveCamera(dx, dy, dz float64) {
	e.camera.X += dx
	e.camera.Y += dy
	e.camera.Z += dz
}

// Update recalculates the world matrix
func (e *Engine) Update() {
	// Always start from the identity matrix
	e.world = Identity4x4()

	// Apply rotations and translations to the world matrix
	e.world = e.world.MulM(&e.rotX)
	e.world = e.world.MulM(&e.rotY)
	e.world = e.world.MulM(&e.rotZ)
	e.world = e.world.MulM(&e.trans)

	up := Vector3d{X: 0, Y: 1, Z: 0}
	target := Vector3d{X: 0, Y: 0, Z: 1}

	// Apply camera rotations, X and Y axes are supported at the moment
	// TODO: figure out why spinning the camera around the Z axis isn't working
	cameraRotation := Identity4x4()
	cameraRotation.RotateY(e.yaw)
	cameraRotation.RotateX(e.pitch)

	e.direction = cameraRotation.MulV(&target)
	target = e.camera.Add(&e.direction)

	cameraMatrix := Identity4x4()
	cameraMatrix.PointAt(&e.camera, &target, &up)
	e.view = cameraMatrix.Inverse()
}

// drawTriangle draw all pixels of a triangle. Supports textured and colored
// triangles
func (e *Engine) drawTriangle(triangle *Triangle, userData UserData) {
	x1, y1, u1, v1, w1 := triangle.UnpackVertex(0)
	x2, y2, u2, v2, w2 := triangle.UnpackVertex(1)
	x3, y3, u3, v3, w3 := triangle.UnpackVertex(2)

	// Presort points depending on the shape of the triangle
	if y2 < y1 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
		w1, w2 = w2, w1
	}

	if y3 < y1 {
		y1, y3 = y3, y1
		x1, x3 = x3, x1
		u1, u3 = u3, u1
		v1, v3 = v3, v1
		w1, w3 = w3, w1
	}

	if y3 < y2 {
		y2, y3 = y3, y2
		x2, x3 = x3, x2
		u2, u3 = u3, u2
		v2, v3 = v3, v2
		w2, w3 = w3, w2
	}

	// Calculate offsets
	dy1 := y2 - y1
	dx1 := x2 - x1
	dv1 := v2 - v1
	du1 := u2 - u1
	dw1 := w2 - w1

	dy2 := y3 - y1
	dx2 := x3 - x1
	dv2 := v3 - v1
	du2 := u3 - u1
	dw2 := w3 - w1

	daxStep := 0.0
	dbxStep := 0.0
	du1Step := 0.0
	dv1Step := 0.0
	dw1Step := 0.0
	du2Step := 0.0
	dv2Step := 0.0
	dw2Step := 0.0

	// Calculate steps
	if dy1 >= 0 {
		daxStep = float64(dx1) / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		dbxStep = float64(dx2) / math.Abs(float64(dy2))
	}

	if dy1 >= 0 {
		du1Step = du1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dv1Step = dv1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dw1Step = dw1 / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		du2Step = du2 / math.Abs(float64(dy2))
	}
	if dy2 >= 0 {
		dv2Step = dv2 / math.Abs(float64(dy2))
	}
	if dy2 >= 0 {
		dw2Step = dw2 / math.Abs(float64(dy2))
	}

	if dy1 >= 0 {
		if e.textureAtlas == nil && triangle.Color != nil {
			for i := y1; i <= y2; i++ {
				ax := float64(x1) + float64(i-y1)*daxStep
				bx := float64(x1) + float64(i-y1)*dbxStep
				texSw := w1 + float64(i-y1)*dw1Step
				texEw := w1 + float64(i-y1)*dw2Step

				if ax > bx {
					ax, bx = bx, ax
					texSw, texEw = texEw, texSw
				}

				texW := texSw
				tStep := 1.0 / (bx - ax)
				t := 0.0

				for j := ax; j < bx; j++ {
					texW = (1.0-t)*texSw + t*texEw

					if texW > e.depthBuffer.At(int(j), i) {
						e.drawPixel(int(j), i, triangle.Color, userData)
						e.depthBuffer.Set(int(j), i, texW)
					}

					t += tStep
				}
			}
		} else if e.textureAtlas != nil {
			for i := y1; i <= y2; i++ {
				ax := float64(x1) + float64(i-y1)*daxStep
				bx := float64(x1) + float64(i-y1)*dbxStep

				texSu := u1 + float64(i-y1)*du1Step
				texSv := v1 + float64(i-y1)*dv1Step
				texSw := w1 + float64(i-y1)*dw1Step

				texEu := u1 + float64(i-y1)*du2Step
				texEv := v1 + float64(i-y1)*dv2Step
				texEw := w1 + float64(i-y1)*dw2Step

				if ax > bx {
					ax, bx = bx, ax
					texSu, texEu = texEu, texSu
					texSv, texEv = texEv, texSv
					texSw, texEw = texEw, texSw
				}

				texU := texSu
				texV := texSv
				texW := texSw

				tStep := 1.0 / (bx - ax)
				t := 0.0

				for j := ax; j < bx; j++ {
					texU = (1.0-t)*texSu + t*texEu
					texV = (1.0-t)*texSv + t*texEv
					texW = (1.0-t)*texSw + t*texEw

					textureWidth := float64(e.textureAtlas.W() - 1)
					textureHeight := float64(e.textureAtlas.H() - 1)

					if texW > e.depthBuffer.At(int(j), i) {
						textureX := int((texU / texW) * textureWidth)
						textureY := 0
						if e.yOrigin == YOriginUpperLeft {
							textureY = int((texV / texW) * textureHeight)
						} else {
							textureY = int((1 - texV/texW) * textureHeight) // invert Y to conform with blender origin
						}
						e.drawPixel(int(j), i, e.textureAtlas.ColorAt(textureX, textureY), userData)
						e.depthBuffer.Set(int(j), i, texW)
					}
					t += tStep
				}
			}
		} else {
			panic("draw error: neither textureAtlas nor color defined")
		}
	}

	dy1 = y3 - y2
	dx1 = x3 - x2
	dv1 = v3 - v2
	du1 = u3 - u2
	dw1 = w3 - w2

	if dy1 >= 0 {
		daxStep = float64(dx1) / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		dbxStep = float64(dx2) / math.Abs(float64(dy2))
	}

	du1Step = 0
	dv1Step = 0

	if dy1 >= 0 {
		du1Step = du1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dv1Step = dv1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dw1Step = dw1 / math.Abs(float64(dy1))
	}

	if dy1 >= 0 {
		if e.textureAtlas == nil && triangle.Color != nil {
			for i := y2; i <= y3; i++ {
				ax := float64(x2) + float64(i-y2)*daxStep
				bx := float64(x1) + float64(i-y1)*dbxStep
				texSw := w2 + float64(i-y2)*dw1Step
				texEw := w1 + float64(i-y1)*dw2Step

				if ax > bx {
					ax, bx = bx, ax
					texSw, texEw = texEw, texSw
				}

				texW := texSw
				tStep := 1.0 / (bx - ax)
				t := 0.0

				for j := ax; j < bx; j++ {
					texW = (1.0-t)*texSw + t*texEw
					if texW > e.depthBuffer.At(int(j), i) {
						e.drawPixel(int(j), i, triangle.Color, userData)
						e.depthBuffer.Set(int(j), i, texW)
					}
					t += tStep
				}
			}

		} else if e.textureAtlas != nil {
			for i := y2; i <= y3; i++ {
				ax := float64(x2) + float64(i-y2)*daxStep
				bx := float64(x1) + float64(i-y1)*dbxStep

				texSu := u2 + float64(i-y2)*du1Step
				texSv := v2 + float64(i-y2)*dv1Step
				texSw := w2 + float64(i-y2)*dw1Step

				texEu := u1 + float64(i-y1)*du2Step
				texEv := v1 + float64(i-y1)*dv2Step
				texEw := w1 + float64(i-y1)*dw2Step

				if ax > bx {
					ax, bx = bx, ax
					texSu, texEu = texEu, texSu
					texSv, texEv = texEv, texSv
					texSw, texEw = texEw, texSw
				}

				texU := texSu
				texV := texSv
				texW := texSw

				tStep := 1.0 / (bx - ax)
				t := 0.0

				for j := ax; j < bx; j++ {
					texU = (1.0-t)*texSu + t*texEu
					texV = (1.0-t)*texSv + t*texEv
					texW = (1.0-t)*texSw + t*texEw

					textureWidth := float64(e.textureAtlas.W() - 1)
					textureHeight := float64(e.textureAtlas.H() - 1)

					if texW > e.depthBuffer.At(int(j), i) {
						textureX := int((texU / texW) * textureWidth)
						textureY := 0
						if e.yOrigin == YOriginUpperLeft {
							textureY = int((texV / texW) * textureHeight)
						} else {
							textureY = int((1 - texV/texW) * textureHeight) // invert Y to conform with blender origin
						}
						e.drawPixel(int(j), i, e.textureAtlas.ColorAt(textureX, textureY), userData)
						e.depthBuffer.Set(int(j), i, texW)
					}

					t += tStep
				}
			}
		} else {
			panic("draw error: neither textureAtlas nor color defined")
		}
	}
}

// renderMesh renders a single mesh
func (e *Engine) renderMesh(mesh Mesh, userData UserData) int {
	// trianglesToRaster holds all visible triangles
	var trianglesToRaster []Triangle
	totalTrianglesRendered := 0

	for _, triangle := range mesh {
		// make copies of all triangle to avoid in place modification
		triangleTransformed := Triangle{}

		// Apply the world matrix
		triangleTransformed.UVs = triangle.UVs.Copy()
		triangleTransformed.Color = triangle.Color
		triangleTransformed.Vertices[0] = e.world.MulV(&triangle.Vertices[0])
		triangleTransformed.Vertices[1] = e.world.MulV(&triangle.Vertices[1])
		triangleTransformed.Vertices[2] = e.world.MulV(&triangle.Vertices[2])

		// Compute the normal for the triangle. They are used to determine if a triangle is visible
		normal := triangleTransformed.Normal()
		cameraRay := triangleTransformed.Vertices[0].Sub(&e.camera)

		// Is the triangle visible?
		dp := normal.Dot(&cameraRay)
		if dp < 0 {
			triangleViewed := Triangle{}

			// Convert world space to view space
			triangleViewed.UVs = triangleTransformed.UVs.Copy()
			triangleViewed.Color = triangleTransformed.Color
			triangleViewed.Vertices[0] = e.view.MulV(&triangleTransformed.Vertices[0])
			triangleViewed.Vertices[1] = e.view.MulV(&triangleTransformed.Vertices[1])
			triangleViewed.Vertices[2] = e.view.MulV(&triangleTransformed.Vertices[2])

			// Check if the triangles are intersecting with screen boundaries and need to be clipped
			p0 := Vector3d{X: 0, Y: 0, Z: 0.1}
			p1 := Vector3d{X: 0, Y: 0, Z: 2.1}
			clippedTriangles := [2]Triangle{}
			numberOfClippedTriangles := triangleViewed.ClipAgainstPlane(&p0, &p1, &clippedTriangles[0], &clippedTriangles[1])

			for n := 0; n < numberOfClippedTriangles; n++ {
				triangleProjected := Triangle{}
				triangleProjected.UVs = clippedTriangles[n].UVs.Copy()

				// Project from 3D into 2D
				triangleProjected.Vertices[0] = e.projection.MulV(&clippedTriangles[n].Vertices[0])
				triangleProjected.Vertices[1] = e.projection.MulV(&clippedTriangles[n].Vertices[1])
				triangleProjected.Vertices[2] = e.projection.MulV(&clippedTriangles[n].Vertices[2])

				triangleProjected.UVs.ScaleW(&triangleProjected)
				triangleProjected.ScaleW()
				triangleProjected.Color = clippedTriangles[n].Color

				// X/Y are inverted so put them back
				triangleProjected.Vertices[0].X *= -1.0
				triangleProjected.Vertices[1].X *= -1.0
				triangleProjected.Vertices[2].X *= -1.0
				triangleProjected.Vertices[0].Y *= -1.0
				triangleProjected.Vertices[1].Y *= -1.0
				triangleProjected.Vertices[2].Y *= -1.0

				offsetView := Vector3d{1, 1, 0, 1}

				triangleProjected.Vertices[0] = triangleProjected.Vertices[0].Add(&offsetView)
				triangleProjected.Vertices[1] = triangleProjected.Vertices[1].Add(&offsetView)
				triangleProjected.Vertices[2] = triangleProjected.Vertices[2].Add(&offsetView)

				triangleProjected.Vertices[0].X *= 0.5 * e.W
				triangleProjected.Vertices[0].Y *= 0.5 * e.H
				triangleProjected.Vertices[1].X *= 0.5 * e.W
				triangleProjected.Vertices[1].Y *= 0.5 * e.H
				triangleProjected.Vertices[2].X *= 0.5 * e.W
				triangleProjected.Vertices[2].Y *= 0.5 * e.H

				// The triangle is ready for rendering
				trianglesToRaster = append(trianglesToRaster, triangleProjected)
			}
		}
	}

	for _, triangle := range trianglesToRaster {
		clipped := [2]Triangle{}
		var finalTrianglesList = []Triangle{triangle}
		newTriangles := 1

		for p := 0; p < 4; p++ {
			trianglesToAdd := 0
			for newTriangles > 0 {
				test := finalTrianglesList[0]
				finalTrianglesList = finalTrianglesList[1:]
				newTriangles--

				switch p {
				case 0:
					trianglesToAdd = test.ClipAgainstPlane(&Vector3d{0, 0, 0, 1}, &Vector3d{0, 1, 0, 1}, &clipped[0], &clipped[1])
					break
				case 1:
					trianglesToAdd = test.ClipAgainstPlane(&Vector3d{0, e.H - 1, 0, 1}, &Vector3d{0, -1, 0, 1}, &clipped[0], &clipped[1])
					break
				case 2:
					trianglesToAdd = test.ClipAgainstPlane(&Vector3d{0, 0, 0, 1}, &Vector3d{1, 0, 0, 1}, &clipped[0], &clipped[1])
					break
				case 3:
					trianglesToAdd = test.ClipAgainstPlane(&Vector3d{e.W - 1, 0, 0, 1}, &Vector3d{-1, 0, 0, 1}, &clipped[0], &clipped[1])
					break
				}

				for w := 0; w < trianglesToAdd; w++ {
					finalTrianglesList = append(finalTrianglesList, clipped[w])
				}
			}
			newTriangles = len(finalTrianglesList)
		}

		for _, t := range finalTrianglesList {
			e.drawTriangle(&t, userData)
			totalTrianglesRendered++
		}
	}

	return totalTrianglesRendered
}

// Render renders all meshes
func (e *Engine) Render(userData UserData) {
	start := time.Now().UnixMilli()

	e.depthBuffer.Clear()

	totalTrianglesRendered := 0
	for _, mesh := range e.meshes {
		totalTrianglesRendered += e.renderMesh(mesh, userData)
	}

	finish := time.Now().UnixMilli()
	e.Metrics.RenderTime = finish - start
	e.Metrics.Triangles = totalTrianglesRendered
}

// NewEngine creates a new 3d engine instance with the given internal
// width and height
func NewEngine(w, h int, fovDegrees float64, drawHook DrawHook, opts *EngineOptions) *Engine {
	engine := &Engine{w: w, h: h, W: float64(w), H: float64(h)}
	engine.meshes = make([]Mesh, 0)
	engine.rotX = Identity4x4()
	engine.rotY = Identity4x4()
	engine.rotZ = Identity4x4()
	engine.trans = Identity4x4()
	engine.camera = Vector3d{
		X: 0,
		Y: 0,
		Z: 0,
		W: 1,
	}

	aspectRatio := float64(w) / float64(h)
	fov := 1.0 / math.Tan(fovDegrees*0.5/180*math.Pi)
	engine.projection = Projection4x4(fov, aspectRatio, 0.1, 1000)
	engine.depthBuffer = NewDepthBuffer(w, h)
	engine.drawPixel = drawHook
	engine.yOrigin = opts.GetYOrigin()
	engine.textureAtlas = opts.GetTextureAtlas()

	return engine
}
