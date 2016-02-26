package main

import (
	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type TestPuzzle3 struct{}

func (c *TestPuzzle3) Preload() {
}

func (c *TestPuzzle3) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})
	w.AddSystem(&TestPuzzle3System{})

	menu := ecs.NewEntity([]string{"RenderSystem", "TestPuzzle3System"})
	w.AddEntity(menu)
}

func (*TestPuzzle3) Hide()        {}
func (*TestPuzzle3) Show()        {}
func (*TestPuzzle3) Type() string { return "TestPuzzle3" }

type TestPuzzle3System struct {
	*ecs.System

	entities []*Entity

	world *ecs.World

	dragging *Entity
	lastpos  engi.Point
}

func (TestPuzzle3System) Type() string {
	return "TestPuzzle3System"
}

func (c *TestPuzzle3System) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	//usuwanie poprzedniej takiej samej sceny, jeśli istniała,  nie przepełniać pamięci, można później przenieść do zamykania po prostu, albo nie pracować na scenach
	if _, ok := entititesFunctions[engi.CurrentScene()]; ok {
		delete(entititesFunctions, engi.CurrentScene())
	}

	entwindow := NewEntity("window", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitestblack.png"),
		Position: engi.Point{X: 0, Y: 0},
		Width:    1024,
		Height:   768,
		Priority: engi.MiddleGround,
	})

	entwindow.AddChildren(NewEntity("windowtopborder", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:                  loadTexture("guitesttop.png"),
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		Width:          1024,
		Height:         20,
		MoveWithParent: true,
		Priority:       engi.MiddleGround + 1,
		OnDragged: func(e *Entity) {
			e.Parent.PosAdd(engi.Point{X: mouseX - prevMouseX, Y: mouseY - prevMouseY})
			minx := CameraGetZ()*(engi.Width()/engi.WindowWidth()) + CameraGetX() - (engi.Width()/2)*CameraGetZ()
			miny := CameraGetZ()*(engi.Height()/engi.WindowHeight()) + CameraGetY() - (engi.Height()/2)*CameraGetZ()
			if e.Parent.Space.Position.X < minx {
				e.Parent.PosSet(engi.Point{X: minx, Y: e.Parent.Space.Position.Y})
			} else if e.Parent.Space.Position.X > minx+engi.Width()-e.Parent.Space.Width {
				e.Parent.PosSet(engi.Point{X: minx + engi.Width() - e.Parent.Space.Width, Y: e.Parent.Space.Position.Y})
			}
			if e.Parent.Space.Position.Y < miny {
				e.Parent.PosSet(engi.Point{X: e.Parent.Space.Position.X, Y: miny})
			} else if e.Parent.Space.Position.Y > miny+engi.Height()-e.Parent.Space.Height {
				e.Parent.PosSet(engi.Point{X: e.Parent.Space.Position.X, Y: miny + engi.Height() - e.Parent.Space.Height})
			}
		},
	}))

	entwindow.AddChildren(NewEntity("windowtopleftborder", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:                  loadTexture("guitestleft.png"),
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		Width:          20,
		Height:         20,
		MoveWithParent: true,
		Priority:       engi.MiddleGround + 2,
	}))

	entwindow.AddChildren(NewEntity("windowtoprightcloseicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:                  loadTexture("guix.png"),
		PositionRelativeToParent: engi.Point{X: 1024 - 20, Y: 0},
		Width:          20,
		Height:         20,
		MoveWithParent: true,
		Priority:       engi.MiddleGround + 2,
		OnPress: func(e *Entity) {
			engi.SetSceneByName("Gui", false)
		},
	}))
}

func (c *TestPuzzle3System) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities(dt)
}
