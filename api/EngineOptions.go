package api

type YOrigin int

const (
	YOriginUpperLeft YOrigin = iota
	YOriginLowerLeft         // Blender exports UV coordinates with the origin in the lower left
)

type EngineOptions struct {
	TextureAtlas TextureAtlas
	YOrigin      YOrigin
}

func (e *EngineOptions) GetYOrigin() YOrigin {
	if e == nil {
		return YOriginUpperLeft
	}
	return e.YOrigin
}

func (e *EngineOptions) GetTextureAtlas() TextureAtlas {
	if e == nil {
		return nil
	}
	return e.TextureAtlas
}
