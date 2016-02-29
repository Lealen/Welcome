package main

import (
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Editor struct{}

func (c *Editor) Preload() {
}

func (c *Editor) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})

	w.AddSystem(&EditorSystem{})

	editor := ecs.NewEntity([]string{"RenderSystem", "EditorSystem"})
	w.AddEntity(editor)
}

func (*Editor) Hide()        {}
func (*Editor) Show()        {}
func (*Editor) Type() string { return "Editor" }

type EditorSystem struct {
	*ecs.System

	smallFont,
	mediumFont,
	bigFont *engi.Font

	world *ecs.World

	dragging *Entity
	lastpos  engi.Point
}

func (EditorSystem) Type() string {
	return "EditorSystem"
}

func (c *EditorSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	menufontname := "04b.ttf"

	c.smallFont = GetFont(menufontname, 12, color.Transparent, color.White)
	c.mediumFont = GetFont(menufontname, 24, color.Transparent, color.White)
	c.bigFont = GetFont(menufontname, 40, color.Transparent, color.White)

	NewEntity("editorhudbg", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 80,
		}),
		Position: engi.Point{X: 0, Y: 0},
		Width:    200,
		Height:   engi.Height(),
		Priority: engi.HUDGround,
		OnWindowResize: func(e *Entity) {
			e.Space.Height = engi.Height()
			e.Render.SetScale(engi.Point{e.Render.Scale().X, engi.Height()})
		},
	})

	NewEntity("testtext", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  c.mediumFont.Render("text string"),
		Position: engi.Point{X: 40, Y: 40},
		Priority: engi.HUDGround + 1,
	})

	NewEntity("testicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("icon.png"),
		Position: engi.Point{X: 1200, Y: 100},
		Scale:    engi.Point{X: 10, Y: 10},
		Width:    160,
		Height:   160,
	})
}

//var changedt float32

func (c *EditorSystem) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities(dt)

	if engi.Keys.Get(engi.D).JustPressed() && engi.Keys.Get(engi.S).Down() {
		engi.SetSceneByName("Dalt", false)
		return
	}

	const scrollSpeed = 500
	if engi.Keys.Get(engi.ArrowUp).Down() {
		CameraMoveY(-scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.ArrowDown).Down() {
		CameraMoveY(scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.ArrowLeft).Down() {
		CameraMoveX(-scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.ArrowRight).Down() {
		CameraMoveX(scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.PageUp).Down() {
		CameraMoveZ(0.005 * -scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.PageDown).Down() {
		CameraMoveZ(0.005 * scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.Home).Down() {
		CameraSetZ(1)
	}
	if engi.Keys.Get(engi.End).Down() {
		CameraSetX(engi.Width() / 2)
		CameraSetY(engi.Height() / 2)
		CameraSetZ(1)
	}
}
