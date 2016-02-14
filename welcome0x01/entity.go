package main

import (
	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Entity struct {
	Entity    *ecs.Entity
	Render    *engi.RenderComponent
	Space     *engi.SpaceComponent
	Animation *engi.AnimationComponent

	World *ecs.World
}

func NewEntity(name string, texture engi.Drawable, world *ecs.World) *Entity {
	return NewEntityWithSystems(name, texture, []string{"RenderSystem"}, world)
}

func NewEntityWithAnimations(name string, spriteSheet *engi.Spritesheet, rate float32, world *ecs.World) (c *Entity) {
	c = NewEntityWithSystems(name, spriteSheet.Cell(0), []string{"AnimationSystem", "RenderSystem"}, world)
	c.Animation = engi.NewAnimationComponent(spriteSheet.Drawables(), rate)
	c.Entity.AddComponent(c.Animation)
	return
}

func NewEntityWithSystems(name string, texture engi.Drawable, systems []string, world *ecs.World) (c *Entity) {
	c = &Entity{}

	c.Entity = ecs.NewEntity(systems)

	c.World = world

	c.Render = engi.NewRenderComponent(texture, engi.Point{X: 1, Y: 1}, name)

	width := texture.Width() * c.Render.Scale().X
	height := texture.Height() * c.Render.Scale().Y

	c.Space = &engi.SpaceComponent{Position: engi.Point{X: (engi.Width() - width) / 2, Y: (engi.Height() - height) / 2}, Width: width, Height: height}

	c.Entity.AddComponent(c.Render)
	c.Entity.AddComponent(c.Space)

	c.World.AddEntity(c.Entity)

	return
}
