package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pb82/mini3d/api"
	"image"
	"image/color"
	"log"

	_ "embed"
	_ "image/png"
)

var (
	//go:embed ryu.png
	textureData []byte

	InternalWidth  = 256
	InternalHeight = 256
)

type TextureAtlasImpl struct {
	w, h int
	img  image.Image
}

func (t *TextureAtlasImpl) W() int {
	return t.w
}

func (t *TextureAtlasImpl) H() int {
	return t.h
}

func (t *TextureAtlasImpl) ColorAt(x, y int) color.Color {
	return t.img.At(x, y)
}

func (t *TextureAtlasImpl) LoadTexture() {
	texture, _, _ := image.Decode(bytes.NewReader(textureData))
	t.img = texture
	t.w = texture.Bounds().Dx()
	t.h = texture.Bounds().Dy()
}

type Game struct {
	Engine *api.Engine
}

func (g *Game) Update() error {
	g.Engine.Update()
	return nil
}

func draw(x, y int, c color.Color, userData api.UserData) {
	screen := userData.(*ebiten.Image)
	screen.Set(x, y, c)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.Engine.Render(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return InternalWidth, InternalHeight
}

func main() {
	mesh := api.Mesh{}
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 0, 0},
			{0, 1, 0},
			{1, 1, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))

	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 0, 0},
			{1, 1, 0},
			{1, 0, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	atlas := &TextureAtlasImpl{}
	atlas.LoadTexture()

	opts := &api.EngineOptions{
		TextureAtlas: atlas,
	}

	engine := api.NewEngine(256, 256, 90, draw, opts)
	engine.AddMesh(mesh)

	game := &Game{engine}

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("3D Engine Demo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
