package main

import (
	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type PlayerSystem struct {
	*ecs.System

	player *Entity

	world *ecs.World
}

func (PlayerSystem) Type() string {
	return "PlayerSystem"
}

func (c *PlayerSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	w.AddSystem(&engi.RenderSystem{})

	c.player = NewEntity("player", loadTexture("heart3.png"), c.world)

	//engi.GetCamera().FollowEntity(c.player.Entity)
}

func (c *PlayerSystem) Update(entity *ecs.Entity, dt float32) {
	if engi.Keys.Get(engi.NumThree).Down() {
		engi.SetSceneByName("Menu", false)
		return
	}

	up := engi.Keys.Get(engi.W).Down() || engi.Keys.Get(engi.ArrowUp).Down()
	down := engi.Keys.Get(engi.S).Down() || engi.Keys.Get(engi.ArrowDown).Down()
	left := engi.Keys.Get(engi.A).Down() || engi.Keys.Get(engi.ArrowLeft).Down()
	right := engi.Keys.Get(engi.D).Down() || engi.Keys.Get(engi.ArrowRight).Down()

	var movementspeed float32 = 200.0

	if up {
		c.player.Space.Position.Y -= movementspeed * dt
	}
	if down {
		c.player.Space.Position.Y += movementspeed * dt
	}
	if left {
		c.player.Space.Position.X -= movementspeed * dt
	}
	if right {
		c.player.Space.Position.X += movementspeed * dt
	}
}
