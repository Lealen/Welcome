package main

import (
	"image"
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type BorderSystem struct {
	*ecs.System

	hidden bool

	borderleft,
	borderright,
	bordertop,
	borderbottom *Entity

	world *ecs.World
}

func (BorderSystem) Type() string {
	return "BorderSystem"
}

func (c *BorderSystem) New(world *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = world

	blackcolor := image.NewGray(image.Rect(0, 0, 1, 1))
	blackcolor.Set(0, 0, color.Black)

	blackimage := blackcolor.SubImage(image.Rect(0, 0, 1, 1))

	blacktexture := engi.NewTexture(engi.LoadImage(blackimage))

	c.borderleft = NewEntity("cameraborderleft", blacktexture, world)
	c.borderleft.Render.SetPriority(39)
	c.borderleft.Render.SetScale(engi.Point{X: 2000, Y: engi.Height()})
	c.borderleft.Space.Position.Set(-2000, 0)
	c.borderleft.Space.Height = engi.Height()
	c.borderleft.Space.Width = 2000

	c.borderright = NewEntity("cameraborderright", blacktexture, world)
	c.borderright.Render.SetPriority(39)
	c.borderright.Render.SetScale(engi.Point{X: 2000, Y: engi.Height()})
	c.borderright.Space.Position.Set(engi.Width(), 0)
	c.borderright.Space.Height = engi.Height()
	c.borderright.Space.Width = 2000

	//TODO: napraw dziwny błąd powodujący nie pokazywanie się jednej z tych dwój barier (left, top)

	/*

		c.bordertop = NewEntity("camerabordertop", blacktexture, world)
		c.bordertop.Render.SetPriority(39)
		c.bordertop.Render.SetScale(engi.Point{X: engi.Width(), Y: 2000})
		c.bordertop.Space.Position.Set(0, -2000)
		c.bordertop.Space.Height = 2000
		c.bordertop.Space.Width = engi.Width()

		c.borderbottom = NewEntity("cameraborderbottom", blacktexture, world)
		c.borderbottom.Render.SetPriority(39)
		c.borderbottom.Render.SetScale(engi.Point{X: engi.Width(), Y: 2000})
		c.borderbottom.Space.Position.Set(0, engi.Height())
		c.borderbottom.Space.Height = 2000
		c.borderbottom.Space.Width = engi.Width()

	   //*/
}

func (c *BorderSystem) Update(entity *ecs.Entity, dt float32) {
	c.borderleft.Space.Position.Set(engi.GetCamera().X()-engi.Width()/2-2000, engi.GetCamera().Y()-engi.Height()/2)
	c.borderright.Space.Position.Set(engi.GetCamera().X()+engi.Width()/2, engi.GetCamera().Y()-engi.Height()/2)
	//TODO: bordertop
	//TODO: borderbottom
}

func (c *BorderSystem) SwitchVisibility() {
	if c.hidden {
		c.borderleft.Render.SetPriority(39)
		c.borderright.Render.SetPriority(39)
	} else {
		c.borderleft.Render.SetPriority(-1)
		c.borderright.Render.SetPriority(-1)
	}
	c.hidden = !c.hidden
}
