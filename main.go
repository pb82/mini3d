package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pb82/mini3d/api"
	"image"
	"image/color"
	"log"
	"time"

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
	Engine       *api.Engine
	milliseconds float64
	elapsedTime  float64
	canvas       []byte
}

func (g *Game) Update() error {
	milliseconds := float64(time.Now().UnixMilli())
	delta := milliseconds - g.milliseconds
	g.milliseconds = milliseconds
	g.elapsedTime += delta

	g.Engine.RotateY(g.elapsedTime / 1000)
	g.Engine.RotateX(g.elapsedTime / 1000)

	g.Engine.Update()
	return nil
}

func draw(x, y int, c color.Color, userData api.UserData) {
	canvas := userData.([]byte)
	// screen.Set(x, y, c)
	r, g, b, a := c.RGBA()
	pos := (y * 256 * 4) + x*4
	canvas[pos] = byte(r >> 8)
	canvas[pos+1] = byte(g >> 8)
	canvas[pos+2] = byte(b >> 8)
	canvas[pos+3] = byte(a >> 8)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.canvas {
		g.canvas[i] = 0
	}
	g.Engine.Render(g.canvas)
	screen.WritePixels(g.canvas)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f FPS, %d tris, %d rt", ebiten.ActualFPS(), g.Engine.Metrics.Triangles, g.Engine.Metrics.RenderTime))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return InternalWidth, InternalHeight
}

func main() {
	mesh := api.Mesh{}

	// Side 1
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

	// Side 2
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 0},
			{1, 1, 0},
			{1, 1, 1},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 1},
		},
		api.Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 3
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 1},
			{1, 1, 1},
			{0, 1, 1},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 1},
			{0, 1, 1},
			{0, 0, 1},
		},
		api.Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 4
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 0, 1},
			{0, 1, 1},
			{0, 1, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 0, 1},
			{0, 1, 0},
			{0, 0, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 5
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 1, 0},
			{0, 1, 1},
			{1, 1, 1},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{0, 1, 0},
			{1, 1, 1},
			{1, 1, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{1, 0},
			{1, 1},
		}))

	// Side 6
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 1},
			{0, 0, 1},
			{0, 0, 0},
		},
		api.Matrix3x2{
			{0, 1},
			{0, 0},
			{1, 0},
		}))
	mesh = append(mesh, api.TexturedTriangleFromMatrix(
		api.Matrix3x3{
			{1, 0, 1},
			{0, 0, 0},
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
	engine.Translate(0, 0, 5)
	engine.SetCameraPosition(0, 0, 3)

	game := &Game{
		Engine:       engine,
		milliseconds: float64(time.Now().UnixMilli()),
		canvas:       make([]byte, 256*256*4),
	}

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("3D Engine Demo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
