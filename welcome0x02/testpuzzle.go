package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type TestPuzzle struct{}

func (c *TestPuzzle) Preload() {
}

func (c *TestPuzzle) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})
	w.AddSystem(&TestPuzzleSystem{})

	menu := ecs.NewEntity([]string{"RenderSystem", "TestPuzzleSystem"})
	w.AddEntity(menu)
}

func (*TestPuzzle) Hide()        {}
func (*TestPuzzle) Show()        {}
func (*TestPuzzle) Type() string { return "TestPuzzle" }

type TestPuzzleSystem struct {
	*ecs.System

	entities []*Entity

	font,
	orangefont,
	greenfont *engi.Font

	probpozostalo int

	world *ecs.World

	dragging *Entity
	lastpos  engi.Point
}

func (TestPuzzleSystem) Type() string {
	return "TestPuzzleSystem"
}

func (c *TestPuzzleSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	//usuwanie poprzedniej takiej samej sceny, jeśli istniała,  nie przepełniać pamięci, można później przenieść do zamykania po prostu, albo nie pracować na scenach
	if _, ok := entititesFunctions[engi.CurrentScene()]; ok {
		delete(entititesFunctions, engi.CurrentScene())
	}

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

	c.orangefont = &engi.Font{
		URL:  "04b.ttf",
		Size: 40,
		BG:   color.Transparent,
		FG: color.RGBA{
			R: 255,
			G: 180,
			B: 0,
			A: 255,
		},
	}
	err = c.orangefont.CreatePreloaded()
	if err != nil {
		fmt.Println(err)
	}

	c.greenfont = &engi.Font{
		URL:  "04b.ttf",
		Size: 40,
		BG:   color.Transparent,
		FG: color.RGBA{
			R: 0,
			G: 185,
			B: 30,
			A: 255,
		},
	}
	err = c.greenfont.CreatePreloaded()
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

	entwindow.AddChildren(NewEntity("windowtopborder", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitesttop.png"),
		Position: engi.Point{X: 100, Y: 100},
		Scale:    engi.Point{X: 40, Y: 1},
		Width:    800,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 1,
		OnDragged: func(e *Entity) {
			e.Parent.PosAdd(engi.Point{X: mouseX - prevMouseX, Y: mouseY - prevMouseY})
			if e.Parent.Space.Position.X < 0 {
				e.Parent.PosSet(engi.Point{X: 0, Y: e.Parent.Space.Position.Y})
			} else if e.Parent.Space.Position.X > engi.Width()-e.Parent.Space.Width {
				e.Parent.PosSet(engi.Point{X: engi.Width() - e.Parent.Space.Width, Y: e.Parent.Space.Position.Y})
			}
			if e.Parent.Space.Position.Y < 0 {
				e.Parent.PosSet(engi.Point{X: e.Parent.Space.Position.X, Y: 0})
			} else if e.Parent.Space.Position.Y > engi.Height()-e.Parent.Space.Height {
				e.Parent.PosSet(engi.Point{X: e.Parent.Space.Position.X, Y: engi.Height() - e.Parent.Space.Height})
			}
		},
	}))

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

	entwindow.AddChildren(NewEntity("windowtoprightcloseicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guix.png"),
		Position: engi.Point{X: 880, Y: 100},
		Scale:    engi.Point{X: 1, Y: 1},
		Width:    20,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 780, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 2,
		OnPress: func(e *Entity) {
			engi.SetSceneByName("Gui", false)
		},
	}))

	for i := 0; i < 4; i++ {
		number1 := NewEntity("number1", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
			Texture:        c.font.Render("0"),
			Position:       engi.Point{X: float32(220 + i*2*80), Y: 300},
			Scale:          engi.Point{X: 3, Y: 3},
			MoveWithParent: true,
			Priority:       engi.MiddleGround + 3,
			OnPress: func(e *Entity) {
				newnumber := 1
				if number, ok := e.Variables["number"]; ok {
					newnumber = number.(int)
					newnumber = (newnumber + 1) % 10
				}
				e.Variables["number"] = newnumber
				e.Render.SetDrawable(c.font.Render(strconv.Itoa(newnumber)))
			},
			OnRightPress: func(e *Entity) {
				newnumber := 1
				if number, ok := e.Variables["number"]; ok {
					newnumber = number.(int)
					newnumber--
					if newnumber < 0 {
						newnumber = 9
					}
				}
				e.Variables["number"] = newnumber
				e.Render.SetDrawable(c.font.Render(strconv.Itoa(newnumber)))
			},
		})
		number1.Variables["number"] = 0
		entwindow.AddChildren(number1)
		c.entities = append(c.entities, number1)
	}

	c.probpozostalo = 4
	probpozostalo := NewEntity("probpozostalo", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: c.font.Render(strconv.Itoa(c.probpozostalo)),
		Scale:   engi.Point{X: 2, Y: 2},
		PositionRelativeToParent: engi.Point{X: 10, Y: 20},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 3,
	})
	entwindow.AddChildren(probpozostalo)

	randomnumber := randoms.Intn(10000)

	entwindow.AddChildren(NewEntity("numbercheck", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture: c.font.Render("check numbers"),
		Scale:   engi.Point{X: 2, Y: 2},
		PositionRelativeToParent: engi.Point{X: 100, Y: 400},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 3,
		OnPress: func(e *Entity) {
			guesed := 0
			for k, v := range c.entities {
				actnumber := v.Variables["number"].(int)
				if actnumber == (randomnumber/int(math.Pow(10, float64(k))))%10 {
					v.Render.SetDrawable(c.greenfont.Render(strconv.Itoa(actnumber)))
					guesed++
				} else {
					for i := 0; i < 4; i++ {
						if actnumber == (randomnumber/int(math.Pow(10, float64(i))))%10 {
							v.Render.SetDrawable(c.orangefont.Render(strconv.Itoa(actnumber)))
						}
					}
				}
			}
			if guesed == 4 {
				entwindow.AddChildren(NewEntity("youwin", []string{"RenderSystem"}, c.world, &EntityDefaults{
					Texture:  loadTexture("youwin.png"),
					Position: engi.Point{X: entwindow.Space.Position.X, Y: entwindow.Space.Position.Y + 20},
					Scale:    engi.Point{X: 1, Y: 1},
					PositionRelativeToParent: engi.Point{X: 0, Y: 0},
					DontAutoSetPosition:      true,
					MoveWithParent:           true,
					Priority:                 engi.MiddleGround + 4,
				}))
			}
			c.probpozostalo--
			probpozostalo.Render.SetDrawable(c.font.Render(strconv.Itoa(c.probpozostalo)))
			if c.probpozostalo < 0 && guesed != 4 {
				entwindow.AddChildren(NewEntity("youlosebg", []string{"RenderSystem"}, c.world, &EntityDefaults{
					Texture:  loadTexture("guitestblack.png"),
					Position: engi.Point{X: entwindow.Space.Position.X, Y: entwindow.Space.Position.Y + 20},
					Scale:    engi.Point{X: 40, Y: 29},
					PositionRelativeToParent: engi.Point{X: 0, Y: 0},
					DontAutoSetPosition:      true,
					MoveWithParent:           true,
					Priority:                 engi.MiddleGround + 4,
				}))
				entwindow.AddChildren(NewEntity("youlosetext", []string{"RenderSystem"}, c.world, &EntityDefaults{
					Texture: c.orangefont.Render("you lose"),
					Scale:   engi.Point{X: 3, Y: 3},
					PositionRelativeToParent: engi.Point{X: 120, Y: 250},
					MoveWithParent:           true,
					Priority:                 engi.MiddleGround + 5,
				}))
			}
		},
	}))
}

func (c *TestPuzzleSystem) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities()
}
