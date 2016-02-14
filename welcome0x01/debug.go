package main

import (
	"fmt"
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type DebugSystem struct {
	*ecs.System

	enabled bool

	debugtextinfo *Entity

	font *engi.Font

	world *ecs.World
}

func (DebugSystem) Type() string {
	return "DebugSystem"
}

func (c *DebugSystem) New(world *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = world

	c.font = &engi.Font{
		URL:  "04b.ttf",
		Size: 12,
		BG:   color.Transparent,
		FG:   color.White,
	}
	err := c.font.CreatePreloaded()
	if err != nil {
		fmt.Println(err)
	}

	c.debugtextinfo = NewEntity("debugtextinfo", c.font.Render("DEBUG MODE"), c.world)
	c.debugtextinfo.Render.SetPriority(-1)
	c.debugtextinfo.Space.Position.Set(engi.Width()-80, 10)
}

func (c *DebugSystem) Update(entity *ecs.Entity, dt float32) {
	if engi.Keys.Get(engi.LeftControl).Down() && engi.Keys.Get(engi.LeftShift).Down() {
		if engi.Keys.Get(engi.NumDivide).Down() {
			c.enabled = true
			c.debugtextinfo.Render.SetPriority(39)
		} else if engi.Keys.Get(engi.NumMultiply).Down() {
			c.enabled = false
			c.debugtextinfo.Render.SetPriority(-1)
		}
	}

	if !c.enabled {
		return
	}

	const scrollSpeed = 200
	if engi.Keys.Get(engi.NumEight).Down() {
		CameraMoveY(-scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.NumTwo).Down() {
		CameraMoveY(scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.NumFour).Down() {
		CameraMoveX(-scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.NumSix).Down() {
		CameraMoveX(scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.NumAdd).Down() {
		CameraMoveZ(0.01 * -scrollSpeed * dt)
	}
	if engi.Keys.Get(engi.NumSubtract).Down() {
		CameraMoveZ(0.01 * scrollSpeed * dt)
	}

	if engi.Keys.Get(engi.NumDecimal).JustPressed() {
		for _, v := range c.world.Systems() {
			if v.Type() == "BorderSystem" {
				v.(*BorderSystem).SwitchVisibility()
			}
		}
	}

	c.debugtextinfo.Space.Position.Set(CameraGetX()+engi.Width()/2-80, CameraGetY()-engi.Height()/2+10)
}
