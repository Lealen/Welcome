package main

import (
	"fmt"
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Menu struct{}

func (c *Menu) Preload() {
	engi.Files.AddFromDir("assets", true)
}

func (c *Menu) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&MenuSystem{})
	w.AddSystem(&BorderSystem{})
	w.AddSystem(&DebugSystem{})

	engi.RegisterScene(&Game{})

	menu := ecs.NewEntity([]string{"RenderSystem", "MenuSystem", "BorderSystem", "DebugSystem"})
	w.AddEntity(menu)
}

func (*Menu) Hide()        {}
func (*Menu) Show()        {}
func (*Menu) Type() string { return "Menu" }

type MenuSystem struct {
	*ecs.System

	entities []*Entity

	font *engi.Font

	world *ecs.World
}

func (MenuSystem) Type() string {
	return "MenuSystem"
}

func (c *MenuSystem) New(w *ecs.World) {
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
	texture := c.font.Render("hello")

	ent := NewEntity("menutext", texture, c.world)
	c.entities = append(c.entities, ent)

	entgif := NewEntity("giftest", loadTexture("kaczek.gif"), c.world)
	entgif.Render.SetScale(engi.Point{X: 0.6, Y: 0.6})
	entgif.Space.Position.Set(50, 50)
	c.entities = append(c.entities, entgif)

	ent2 := NewEntity("giftest", loadTexture("kw17ljv.jpg"), c.world)
	ent2.Render.SetScale(engi.Point{X: 0.1, Y: 0.1})
	ent2.Space.Position.SetPoint(engi.Point{X: 300, Y: 300})
	c.entities = append(c.entities, ent2)

	entmagus := NewEntityWithAnimations("magus", engi.NewSpritesheetFromFile("magus.png", 40, 40), 0.12, c.world)
	entmagus.Render.SetPriority(10)
	entmagus.Render.SetScale(engi.Point{X: 4, Y: 4})
	entmagus.Space.Position.SetPoint(engi.Point{X: 400, Y: 200})
	entmagus.Animation.AddAnimationAction(&engi.AnimationAction{Name: "walkdown", Frames: []int{22, 23, 24, 25, 26, 27, 26, 25, 24, 23}})
	entmagus.Animation.SelectAnimationByName("walkdown")
	c.entities = append(c.entities, entmagus)

	engi.WorldBounds.Max.Add(engi.Point{X: 1000, Y: 1000})
	engi.WorldBounds.Min.Subtract(engi.Point{X: 1000, Y: 1000})
}

func (c *MenuSystem) Update(entity *ecs.Entity, dt float32) {
	fmt.Printf("%#v                   \r", engi.Mouse)

	if engi.Keys.Get(engi.NumOne).Down() {
		engi.SetSceneByName("Game", false)
		return
	}

	for i := engi.A; i <= engi.Z; i++ {
		if engi.Keys.Get(i).Down() {
			texture := c.font.Render(string([]byte{byte(i)}))

			c.entities[0].Render.SetDrawable(texture)
			break
		}
	}

	if c.entities[0].Space.Position.X > 640 {
		c.entities[0].Space.Position.X = 0
	} else {
		c.entities[0].Space.Position.Add(engi.Point{X: 200 * dt, Y: 0})
	}
}
