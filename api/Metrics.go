package api

type Metrics struct {
	// Triangles rendered in the last call to `Render`
	Triangles int

	// RenderTime represents the time in milliseconds it took to render all meshes during a single
	// call to `Render`
	RenderTime int64
}
