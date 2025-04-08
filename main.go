package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pb82/mini3d/api"
	"image"
	"image/color"
	"log"
	"time"

	_ "embed"
	_ "image/jpeg"
	_ "image/png"
)

var (
	//go:embed texture.jpg
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
	meshes       []*api.Mesh
}

func (g *Game) Update() error {
	milliseconds := float64(time.Now().UnixMilli())
	delta := milliseconds - g.milliseconds
	g.milliseconds = milliseconds
	g.elapsedTime += delta

	// g.meshes[0].Move(0, 0, 1*delta/1000)
	// p := g.meshes[0].GetPosition()
	// g.meshes[0].MoveRelative(0, 0, -1*delta/5000)

	// c := g.meshes[0].GetCenter()
	// g.meshes[0].RotateXAround(1*g.elapsedTime/1000, &c)

	// g.meshes[1].RotateWorldAroundY(1*g.elapsedTime/1000, -g.meshes[1].GetPosition().X-.5, -g.meshes[1].GetPosition().Z-.5)
	// g.meshes[2].RotateWorldAroundZ(1*g.elapsedTime/1000, -g.meshes[2].GetPosition().X-.5, -g.meshes[2].GetPosition().Y-.5)

	keys := inpututil.AppendPressedKeys([]ebiten.Key{ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight})
	for _, key := range keys {
		if key == ebiten.KeyW {
			g.Engine.MoveCameraForward(1 * delta / 500)
		}
		if key == ebiten.KeyS {
			g.Engine.MoveCameraForward(-1 * delta / 500)
		}
		if key == ebiten.KeyA {
			g.Engine.SetCameraPositionRelative(0, 0, 0, -1*delta/1000, 0)
		}
		if key == ebiten.KeyD {
			g.Engine.SetCameraPositionRelative(0, 0, 0, 1*delta/1000, 0)
		}
		if key == ebiten.KeyUp {
			g.Engine.SetCameraPositionRelative(0, 0, 0, 0, -1*delta/1000)
		}
		if key == ebiten.KeyDown {
			g.Engine.SetCameraPositionRelative(0, 0, 0, 0, 1*delta/1000)
		}
	}

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
	mesh1, err := api.LoadWavefrontObj("./cube.obj")
	if err != nil {
		panic(err)
	}

	mesh2 := mesh1.Copy()
	mesh3 := mesh1.Copy()
	mesh4 := mesh1.Copy()
	mesh5 := mesh1.Copy()
	mesh6 := mesh1.Copy()
	mesh7 := mesh1.Copy()
	mesh8 := mesh1.Copy()
	mesh9 := mesh1.Copy()
	mesh10 := mesh1.Copy()
	mesh11 := mesh1.Copy()
	mesh12 := mesh1.Copy()
	mesh13 := mesh1.Copy()

	atlas := &TextureAtlasImpl{}
	atlas.LoadTexture()

	opts := &api.EngineOptions{
		TextureAtlas: atlas,
	}

	mesh1.Translate(0, 0, 0)
	mesh2.Translate(3, 0, 0)
	mesh3.Translate(0, 0, 2)
	mesh4.Translate(3, 0, 2)
	mesh5.Translate(0, 0, 4)
	mesh6.Translate(3, 0, 4)
	mesh7.Translate(0, 0, 6)
	mesh8.Translate(3, 0, 6)

	mesh9.Translate(2, 0, 6)
	mesh10.Translate(2, 2, 6)
	mesh11.Translate(2, 4, 6)
	mesh12.Translate(2, 6, 6)
	mesh13.Translate(2, 8, 6)

	engine := api.NewEngine(256, 256, 90, draw, opts)
	engine.AddMesh(mesh1)
	engine.AddMesh(mesh2)
	engine.AddMesh(mesh3)
	engine.AddMesh(mesh4)
	engine.AddMesh(mesh5)
	engine.AddMesh(mesh6)
	engine.AddMesh(mesh7)
	engine.AddMesh(mesh8)
	engine.AddMesh(mesh9)
	engine.AddMesh(mesh10)
	engine.AddMesh(mesh11)
	engine.AddMesh(mesh12)
	engine.AddMesh(mesh13)

	engine.SetCameraPositionAbsolute(0, 0, -5, 0, 0)

	game := &Game{
		Engine:       engine,
		milliseconds: float64(time.Now().UnixMilli()),
		canvas:       make([]byte, 256*256*4),
		meshes:       []*api.Mesh{mesh1, mesh2},
	}

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("3D Engine Demo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
