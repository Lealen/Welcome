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
	Mouse     *engi.MouseComponent

	Childrens []*Entity
	Parent    *Entity

	PositionRelativeToParent engi.Point
	MoveWithParent,
	DontAutoSetPosition bool

	World *ecs.World

	Initialized,
	IsClicked,
	IsRightClicked bool
	OnFirstUpdate,
	OnClicked,
	OnPress,
	OnRelease,
	OnDragged,
	OnRightClicked,
	OnRightPress,
	OnRightRelease,
	OnEnter,
	OnLeave,
	OnWindowResize,
	OnUpdate func(*Entity)

	Variables map[string]interface{}
}

type EntityDefaults struct {
	Texture     engi.Drawable
	SpriteSheet *engi.Spritesheet

	Position engi.Point
	Scale    engi.Point

	Width,
	Height float32

	AnimationRate float32

	PositionRelativeToParent engi.Point
	MoveWithParent,
	DontAutoSetPosition bool

	Priority engi.PriorityLevel

	OnFirstUpdate,
	OnClicked,
	OnPress,
	OnRelease,
	OnDragged,
	OnRightClicked,
	OnRightPress,
	OnRightRelease,
	OnEnter,
	OnLeave,
	OnWindowResize,
	OnUpdate func(*Entity)
}

func getDefaultDefaults() *EntityDefaults {
	d := EntityDefaults{}
	d.Texture = loadTexture("notfound.png")
	d.SpriteSheet = nil
	d.Position = engi.Point{X: 0, Y: 0}
	d.Scale = engi.Point{X: 1, Y: 1}
	d.Width = 0.0
	d.Height = 0.0
	d.PositionRelativeToParent = engi.Point{X: 0, Y: 0}
	d.MoveWithParent = false
	d.DontAutoSetPosition = false
	d.Priority = engi.MiddleGround
	return &d
}

func NewEntity(name string, systems []string, world *ecs.World, defaults *EntityDefaults) (c *Entity) {
	c = &Entity{}
	c.Entity = ecs.NewEntity(systems)
	entititesFunctions[engi.CurrentScene()] = append(entititesFunctions[engi.CurrentScene()], c)
	c.World = world

	c.Variables = make(map[string]interface{})

	if defaults == nil {
		defaults = getDefaultDefaults()
	}

	c.PositionRelativeToParent = defaults.PositionRelativeToParent
	c.MoveWithParent = defaults.MoveWithParent
	c.DontAutoSetPosition = defaults.DontAutoSetPosition

	for _, v := range systems {
		switch v {
		case "RenderSystem":
			c.Render = engi.NewRenderComponent(defaults.Texture, defaults.Scale, name)
			if defaults.Width == 0 && defaults.Height == 0 {
				defaults.Width = defaults.Texture.Width() * c.Render.Scale().X
				defaults.Height = defaults.Texture.Height() * c.Render.Scale().Y
			}
			c.Entity.AddComponent(c.Render)
		case "AnimationSystem":
			c.Animation = engi.NewAnimationComponent(defaults.SpriteSheet.Drawables(), defaults.AnimationRate)
			c.Entity.AddComponent(c.Animation)
		case "MouseSystem":
			c.Mouse = &engi.MouseComponent{}
			c.Entity.AddComponent(c.Mouse)
		}
	}

	c.Space = &engi.SpaceComponent{Position: defaults.Position, Width: defaults.Width, Height: defaults.Height}
	c.Entity.AddComponent(c.Space)

	c.World.AddEntity(c.Entity)

	if defaults.Priority != engi.MiddleGround {
		c.Render.SetPriority(defaults.Priority)
	}

	c.OnClicked = defaults.OnClicked
	c.OnPress = defaults.OnPress
	c.OnRelease = defaults.OnRelease
	c.OnDragged = defaults.OnDragged
	c.OnRightClicked = defaults.OnRightClicked
	c.OnRightPress = defaults.OnRightPress
	c.OnRightRelease = defaults.OnRightRelease
	c.OnEnter = defaults.OnEnter
	c.OnLeave = defaults.OnLeave
	c.OnWindowResize = defaults.OnWindowResize
	c.OnUpdate = defaults.OnUpdate

	return
}

func (e *Entity) RemoveEntity() {
	for _, v := range e.Childrens {
		v.RemoveEntity()
	}
	for k, v := range entititesFunctions[engi.CurrentScene()] {
		if v == e {
			entititesFunctions[engi.CurrentScene()] = append(entititesFunctions[engi.CurrentScene()][:k], entititesFunctions[engi.CurrentScene()][k+1:]...)
			return
		}
	}
	if e.Render != nil {
		e.Render.SetDrawable(nil)
		e.Entity.RemoveComponent(e.Render)
	}
	if e.Space != nil {
		e.Entity.RemoveComponent(e.Space)
	}
	if e.Animation != nil {
		e.Entity.RemoveComponent(e.Animation)
	}
	if e.Mouse != nil {
		e.Entity.RemoveComponent(e.Mouse)
	}
	e.World.RemoveEntity(e.Entity)
}

func (p *Entity) AddChildren(c *Entity) {
	p.Childrens = append(p.Childrens, c)
	c.Parent = p

	if !c.DontAutoSetPosition && c.PositionRelativeToParent.X == 0 && c.PositionRelativeToParent.Y == 0 {
		c.PositionRelativeToParent = engi.Point{X: c.Space.Position.X - p.Space.Position.X, Y: c.Space.Position.Y - p.Space.Position.Y}
	} else if !c.DontAutoSetPosition && c.Space.Position.X == 0 && c.Space.Position.Y == 0 {
		c.Space.Position.X = p.Space.Position.X + c.PositionRelativeToParent.X
		c.Space.Position.Y = p.Space.Position.Y + c.PositionRelativeToParent.Y
	}
}

func (c *Entity) PosAdd(p2 engi.Point) {
	c.Space.Position.Add(p2)
	if c.Parent != nil && c.MoveWithParent {
		c.PositionRelativeToParent = engi.Point{X: c.Space.Position.X - c.Parent.Space.Position.X, Y: c.Space.Position.Y - c.Parent.Space.Position.Y}
	}
	for _, v := range c.Childrens {
		if v.MoveWithParent {
			v.PosAdd(p2)
		}
	}
}

func (c *Entity) PosSet(p2 engi.Point) {
	c.Space.Position.SetPoint(p2)
	if c.Parent != nil && c.MoveWithParent {
		c.PositionRelativeToParent = engi.Point{X: c.Space.Position.X - c.Parent.Space.Position.X, Y: c.Space.Position.Y - c.Parent.Space.Position.Y}
	}
	for _, v := range c.Childrens {
		if v.MoveWithParent {
			ptmp := p2
			ptmp.Add(v.PositionRelativeToParent)
			v.PosSet(ptmp)
		}
	}
}
