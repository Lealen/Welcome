package main

import (
	"fmt"
	"image/color"
	"os/user"

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
	engi.RegisterScene(&TestPuzzle2{})
	engi.RegisterScene(&TestPuzzle3{})

	engi.RegisterScene(&Editor{})

	engi.RegisterScene(&Dalt{})

	menu := ecs.NewEntity([]string{"RenderSystem", "GuiSystem"})
	w.AddEntity(menu)
}

func (*Gui) Hide()        {}
func (*Gui) Show()        {}
func (*Gui) Type() string { return "Gui" }

type GuiSystem struct {
	*ecs.System

	entitiestext []*Entity

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

	usr, _ := user.Current()
	NewEntity("gotopuzzle", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  c.font.Render(fmt.Sprintf("Hello, %s!", usr.Username)),
		Position: engi.Point{X: 10, Y: 10},
		Scale:    engi.Point{X: 1, Y: 1},
		Priority: engi.HUDGround,
	})

	NewEntity("testicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("icon.png"),
		Position: engi.Point{X: 1200, Y: 100},
		Scale:    engi.Point{X: 10, Y: 10},
		Width:    160,
		Height:   160,
		OnDragged: func(e *Entity) {
			e.PosAdd(engi.Point{X: mouseX - prevMouseX, Y: mouseY - prevMouseY})
			minx := CameraGetZ()*(engi.Width()/engi.WindowWidth()) + CameraGetX() - (engi.Width()/2)*CameraGetZ()
			miny := CameraGetZ()*(engi.Height()/engi.WindowHeight()) + CameraGetY() - (engi.Height()/2)*CameraGetZ()
			if e.Space.Position.X < minx {
				e.PosSet(engi.Point{X: minx, Y: e.Space.Position.Y})
			} else if e.Space.Position.X > minx+engi.Width()-e.Space.Width {
				e.PosSet(engi.Point{X: minx + engi.Width() - e.Space.Width, Y: e.Space.Position.Y})
			}
			if e.Space.Position.Y < miny {
				e.PosSet(engi.Point{X: e.Space.Position.X, Y: miny})
			} else if e.Space.Position.Y > miny+engi.Height()-e.Space.Height {
				e.PosSet(engi.Point{X: e.Space.Position.X, Y: miny + engi.Height() - e.Space.Height})
			}
		},
	})

	NewEntity("gotopuzzle", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  c.font.Render("_test_puzzle_1_"),
		Position: engi.Point{X: engi.Width() - 320, Y: engi.Height() - 50},
		Scale:    engi.Point{X: 1, Y: 1},
		OnPress: func(e *Entity) {
			engi.SetSceneByName("TestPuzzle", true)
		},
	})

	NewEntity("gotopuzzle", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  c.font.Render("_test_puzzle_2_"),
		Position: engi.Point{X: engi.Width() - 320, Y: engi.Height() - 50 - 50},
		Scale:    engi.Point{X: 1, Y: 1},
		OnPress: func(e *Entity) {
			engi.SetSceneByName("TestPuzzle2", true)
		},
	})

	NewEntity("gotopuzzle", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  c.font.Render("_test_puzzle_3_"),
		Position: engi.Point{X: engi.Width() - 320, Y: engi.Height() - 50 - 100},
		Scale:    engi.Point{X: 1, Y: 1},
		OnPress: func(e *Entity) {
			engi.SetSceneByName("TestPuzzle3", true)
		},
	})

}

//var changedt float32

func (c *GuiSystem) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities(dt)

	//if engi.Keys.Get(engi.CONTROL).Down() {
	if engi.Keys.Get(engi.E).JustPressed() {
		err := engi.SetSceneByName("Editor", false)
		if err != nil {
			fmt.Println(err)
		}
	}
	//}
	/*
		changedt += dt
		if changedt > 0.04 {
			for _, v := range c.entitiestext {
				v.PosSetRel(engi.Point{X: randoms.Float32()*560 + 20, Y: randoms.Float32()*540 + 20})
			}
			changedt -= 0.04
		}
		//*/
}
