package main

import (
	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type BoobTities struct {
	ent     *ecs.Entity
	texture *engi.Texture
	render  *engi.RenderComponent
	space   *engi.SpaceComponent
}

type MapSystem struct {
	*ecs.System
	world    *ecs.World
	entities []*BoobTities
}

func (MapSystem) Type() string {
	return "MapSystem"
}

func (c *MapSystem) Load(name string) {
	for _, v := range c.entities {
		c.world.RemoveEntity(v.ent)
	}
}

func (c *MapSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	w.AddSystem(&engi.RenderSystem{})
}

func (c *MapSystem) Update(entity *ecs.Entity, dt float32) {
}
