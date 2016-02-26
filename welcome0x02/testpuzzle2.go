package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"strconv"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
	"github.com/disintegration/imaging"
)

type TestPuzzle2 struct{}

func (c *TestPuzzle2) Preload() {
}

func (c *TestPuzzle2) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})
	w.AddSystem(&TestPuzzle2System{})

	menu := ecs.NewEntity([]string{"RenderSystem", "TestPuzzle2System"})
	w.AddEntity(menu)
}

func (*TestPuzzle2) Hide()        {}
func (*TestPuzzle2) Show()        {}
func (*TestPuzzle2) Type() string { return "TestPuzzle2" }

type TestPuzzle2System struct {
	*ecs.System

	entities []*Entity

	world *ecs.World

	dragging *Entity
	lastpos  engi.Point
}

func (TestPuzzle2System) Type() string {
	return "TestPuzzle2System"
}

func (c *TestPuzzle2System) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	//usuwanie poprzedniej takiej samej sceny, jeśli istniała,  nie przepełniać pamięci, można później przenieść do zamykania po prostu, albo nie pracować na scenach
	if _, ok := entititesFunctions[engi.CurrentScene()]; ok {
		delete(entititesFunctions, engi.CurrentScene())
	}

	entwindow := NewEntity("window", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitestblack.png"),
		Position: engi.Point{X: 100, Y: 100},
		Width:    800,
		Height:   600 + 20,
		Priority: engi.MiddleGround,
	})

	entwindow.AddChildren(NewEntity("windowtopborder", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guitesttop.png"),
		Position: engi.Point{X: 100, Y: 100},
		Width:    800,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 1,
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
		Texture:  loadTexture("guitestleft.png"),
		Position: engi.Point{X: 100, Y: 100},
		Width:    20,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 0, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 2,
	}))

	entwindow.AddChildren(NewEntity("windowtoprightcloseicon", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture:  loadTexture("guix.png"),
		Position: engi.Point{X: 880, Y: 100},
		Width:    20,
		Height:   20,
		PositionRelativeToParent: engi.Point{X: 780, Y: 0},
		MoveWithParent:           true,
		Priority:                 engi.MiddleGround + 2,
		OnPress: func(e *Entity) {
			engi.SetSceneByName("Gui", false)
		},
	}))

	wezFragmentObrazka := func(imgo image.Image, x0, y0, x1, y1 int) *engi.ImageObject {
		img := imaging.Crop(imgo, image.Rect(x0, y0, x1, y1))

		b := img.Bounds()
		newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

		return engi.NewImageObject(newm)
	}

	randomnumber := randoms.Intn(9) + 1 //1-9

	file, err := os.Open("assets" + string(os.PathSeparator) + "testpuzzle2" + string(os.PathSeparator) + strconv.Itoa(randomnumber) + ".jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		return
	}

	var siatka [16]*Entity
	var draggin *Entity
	var przeniesiono bool

	for i := 0; i < 16; i++ {
		piece := NewEntity("number1", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
			Texture: engi.NewTexture(wezFragmentObrazka(img, 800/4*(i%4), 600/4*(i/4), 800/4*(i%4)+800/4, 600/4*(i/4)+600/4)),
			Scale:   engi.Point{X: 1, Y: 1},
			PositionRelativeToParent: engi.Point{X: float32(800 / 4 * (i % 4)), Y: float32(600/4*(i/4)) + 20},
			MoveWithParent:           true,
			Priority:                 engi.MiddleGround + 3,
			OnPress: func(v *Entity) {
				if przeniesiono {
					przeniesiono = false
				} else if draggin == nil {
					draggin = v
					v.Render.SetScale(engi.Point{X: 0.9, Y: 0.9})
					v.Space.Position.Add(engi.Point{X: 10, Y: 7.5})
				} else {
					var rn1, rn2 int
					draggin.Render.SetScale(engi.Point{X: 1, Y: 1})
					draggin.Space.Position.Subtract(engi.Point{X: 10, Y: 7.5})
					for k, v2 := range siatka {
						if v2 == v {
							rn1 = k
						}
						if v2 == draggin {
							rn2 = k
						}
					}
					if rn1 != rn2 {
						siatka[rn1], siatka[rn2] = siatka[rn2], siatka[rn1]
						pos1 := siatka[rn1].Space.Position
						siatka[rn1].PosSet(siatka[rn2].Space.Position)
						siatka[rn2].PosSet(pos1)
						przeniesiono = true
					}
					draggin = nil

					allgood := true
					for k, v := range siatka {
						if k != v.Variables["i"].(int) {
							allgood = false
							break
						}
					}
					if allgood {
						fmt.Println("you win")

					}
				}
			},
		})
		piece.Variables["i"] = i
		entwindow.AddChildren(piece)
		siatka[i] = piece
	}

	for i := 0; i < 50; i++ {
		rn1, rn2 := randoms.Intn(16), randoms.Intn(16)
		siatka[rn1], siatka[rn2] = siatka[rn2], siatka[rn1]
		pos1 := siatka[rn1].Space.Position
		siatka[rn1].PosSet(siatka[rn2].Space.Position)
		siatka[rn2].PosSet(pos1)
	}
}

func (c *TestPuzzle2System) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities(dt)
}
