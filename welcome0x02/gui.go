package main

import (
	"fmt"
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Gui struct{}

func (c *Gui) Preload() {
	engi.Files.AddFromDir("assets", true)
}

func (c *Gui) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})
	w.AddSystem(&GuiSystem{})

	engi.RegisterScene(&TestPuzzle{})

	menu := ecs.NewEntity([]string{"RenderSystem", "GuiSystem"})
	w.AddEntity(menu)
}

func (*Gui) Hide()        {}
func (*Gui) Show()        {}
func (*Gui) Type() string { return "Gui" }

type GuiSystem struct {
	*ecs.System

	//entities []*Entity

	font *engi.Font

	world *ecs.World

	dragging *Entity
	lastpos  engi.Point
}

func (GuiSystem) Type() string {
	return "GuiSystem"
}

func (c *GuiSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	c.font = &engi.Font{
		URL:  "04b.ttf",
		Size: 40,
		BG:   color.Transparent,
		FG:   color.White,
	}
	err := c.font.CreatePreloaded()
	if err != nil {
		fmt.Println(err)
	}

	entwindow := NewEntity("window", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitestblack.png"),
		Position: engi.Point{X: 100, Y: 100},
		Scale:    engi.Point{X: 40, Y: 30},
		Width:    800,
		Height:   600,
		Priority: engi.MiddleGround,
	})
	//c.entities = append(c.entities, entwindow)

	entwindow.AddChildren(NewEntity("windowtopborder", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitesttop.png"),
		Position: engi.Point{X: 100, Y: 100},
		Scale:    engi.Point{X: 40, Y: 1},
		Width:    800,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 1,
		OnClicked: func(e *Entity) {
			//fmt.Println("clicked")
		},
		OnPress: func(e *Entity) {
			fmt.Println("pressed")
		},
		OnRelease: func(e *Entity) {
			fmt.Println("released")
		},
		OnDragged: func(e *Entity) {
			//fmt.Println("dragged")
			e.Parent.PosAdd(engi.Point{X: engi.Mouse.X - PreviousMousePosX, Y: engi.Mouse.Y - PreviousMousePosY})
		},
		OnEnter: func(e *Entity) {
			fmt.Println("entered")
		},
		OnLeave: func(e *Entity) {
			fmt.Println("leaved")
		},
		OnWindowResize: func(e *Entity) {
			fmt.Println("resized")
		},
	}))
	//entwindow.AddChildren(entwindowtopborder)

	entwindow.AddChildren(NewEntity("windowtopleftborder", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitestleft.png"),
		Position: engi.Point{X: 100, Y: 100},
		Scale:    engi.Point{X: 1, Y: 1},
		Width:    20,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 2,
	}))
	//entwindow.AddChildren(entwindowtopleftborder)

	entwindowtoprightcloseicon := NewEntity("windowtoprightcloseicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guix.png"),
		Position: engi.Point{X: 880, Y: 100},
		Scale:    engi.Point{X: 1, Y: 1},
		Width:    20,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 780, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 2,
		OnPress: func(e *Entity) {
			e.Parent.RemoveEntity()
		},
	})
	entwindow.AddChildren(entwindowtoprightcloseicon)

	NewEntity("testicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("icon.png"),
		Position: engi.Point{X: 1200, Y: 100},
		Scale:    engi.Point{X: 10, Y: 10},
		Width:    160,
		Height:   160,
		OnDragged: func(e *Entity) {
			e.PosAdd(engi.Point{X: engi.Mouse.X - PreviousMousePosX, Y: engi.Mouse.Y - PreviousMousePosY})
		},
	})

	//*
	for i := 0; i < 300; i++ {
		entwindowtexttest := NewEntity("windowtexttest", []string{"RenderSystem"}, c.world, &EntityDefaults{
			Texture:  c.font.Render("OMG, hello!"),
			Position: engi.Point{X: float32((i*50)%600 + 15 + 100), Y: float32(((i+30)/30)*50%700 + 100)},
			Scale:    engi.Point{X: 1, Y: 1},
			PositionRelativeToParent: engi.Point{X: float32((i*50)%600 + 15), Y: float32(((i + 30) / 30) % 600)},
			MoveWithParent:           true,
			Priority:                 engi.MiddleGround + 1,
		})
		entwindow.AddChildren(entwindowtexttest)
	}
	//*/

	NewEntity("gotopuzzle", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  c.font.Render("_test_puzzle_"),
		Position: engi.Point{X: 1520, Y: 980},
		Scale:    engi.Point{X: 1, Y: 1},
		OnPress: func(e *Entity) {
			engi.SetSceneByName("TestPuzzle", true)
		},
	})

}

func (c *GuiSystem) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities()
}
