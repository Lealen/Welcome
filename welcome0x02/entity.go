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
	PriorityRelativeToParent engi.PriorityLevel
	MoveWithParent,
	DontAutoSetPosition,
	DontAutoSetPriority,
	IgnoreWhatIsOnTop bool

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
	OnWindowResize func(*Entity)
	OnUpdate func(*Entity, float32)

	Variables map[string]interface{}

	priority engi.PriorityLevel
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
	PriorityRelativeToParent engi.PriorityLevel
	MoveWithParent,
	DontAutoSetPosition,
	DontAutoSetPriority,
	IgnoreWhatIsOnTop bool

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
	OnWindowResize func(*Entity)
	OnUpdate func(*Entity, float32)
}

func NewEntity(name string, systems []string, world *ecs.World, defaults *EntityDefaults) (c *Entity) {
	c = &Entity{}
	c.Entity = ecs.NewEntity(systems)
	entititesFunctions[engi.CurrentScene()] = append(entititesFunctions[engi.CurrentScene()], c)
	c.World = world

	c.Variables = make(map[string]interface{})

	if defaults == nil {
		defaults = &EntityDefaults{}
		defaults.Texture = loadTexture("notfound.png")
	}

	c.PositionRelativeToParent = defaults.PositionRelativeToParent
	c.PriorityRelativeToParent = defaults.PriorityRelativeToParent
	c.MoveWithParent = defaults.MoveWithParent
	c.DontAutoSetPosition = defaults.DontAutoSetPosition
	c.DontAutoSetPriority = defaults.DontAutoSetPriority
	c.IgnoreWhatIsOnTop = defaults.IgnoreWhatIsOnTop

	if defaults.Scale.X == 0 && defaults.Scale.Y == 0 && defaults.Width == 0 && defaults.Height == 0 && defaults.Texture != nil {
		defaults.Scale.X = 1
		defaults.Scale.Y = 1
		defaults.Width = defaults.Texture.Width()
		defaults.Height = defaults.Texture.Height()
	} else if defaults.Scale.X == 0 && defaults.Scale.Y == 0 && defaults.Texture != nil {
		defaults.Scale.X = defaults.Width / defaults.Texture.Width()
		defaults.Scale.Y = defaults.Height / defaults.Texture.Height()
	} else if defaults.Width == 0 && defaults.Height == 0 && defaults.Texture != nil {
		defaults.Width = defaults.Texture.Width() * defaults.Scale.X
		defaults.Height = defaults.Texture.Height() * defaults.Scale.Y
	}

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

	c.priority = defaults.Priority

	entitiesChange = true
	return
}

func (e *Entity) RemoveEntity() {
	for _, v := range e.Childrens {
		v.RemoveEntity()
	}
	for k, v := range entititesFunctions[engi.CurrentScene()] {
		if v == e {
			entititesFunctions[engi.CurrentScene()] = append(entititesFunctions[engi.CurrentScene()][:k], entititesFunctions[engi.CurrentScene()][k+1:]...)
			break
		}
	}
	if e.Render != nil {
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

	entitiesChange = true
}

func (p *Entity) AddChildren(c *Entity) {
	p.Childrens = append(p.Childrens, c)
	c.Parent = p

	if !c.DontAutoSetPosition && c.Space.Position.X == 0 && c.Space.Position.Y == 0 {
		c.Space.Position.X = p.Space.Position.X + c.PositionRelativeToParent.X
		c.Space.Position.Y = p.Space.Position.Y + c.PositionRelativeToParent.Y
	} else if !c.DontAutoSetPosition && c.PositionRelativeToParent.X == 0 && c.PositionRelativeToParent.Y == 0 {
		c.PositionRelativeToParent = engi.Point{X: c.Space.Position.X - p.Space.Position.X, Y: c.Space.Position.Y - p.Space.Position.Y}
	}

	if !c.DontAutoSetPriority && c.PriorityRelativeToParent == 0 {
		c.PriorityRelativeToParent = c.priority - p.priority
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

func (c *Entity) PosSetRel(p2 engi.Point) {
	if c.Parent != nil {
		c.Space.Position.SetPoint(c.Parent.Space.Position)
		c.Space.Position.Add(p2)
		c.PositionRelativeToParent = engi.Point{X: c.Space.Position.X - c.Parent.Space.Position.X, Y: c.Space.Position.Y - c.Parent.Space.Position.Y}
	} else {
		c.Space.Position.SetPoint(p2)
	}
	for _, v := range c.Childrens {
		if v.MoveWithParent {
			ptmp := c.Space.Position
			ptmp.Add(v.PositionRelativeToParent)
			v.PosSet(ptmp)
		}
	}
}

func (c *Entity) GetPriority() engi.PriorityLevel {
	return c.priority
}

func (c *Entity) SetPriority(p engi.PriorityLevel) {
	c.priority = p
	c.Render.SetPriority(p)
	entitiesChange = true

	for _, v := range c.Childrens {
		v.SetPriority(p + v.PriorityRelativeToParent)
	}

	if !c.DontAutoSetPriority && c.Parent != nil {
		c.PriorityRelativeToParent = c.priority - c.Parent.priority
	}
}

func AddOneBelowAnother(position engi.Point, shift float32, entities ...*Entity) {
	lastpos := position
	for _, v := range entities {
		v.PosSet(lastpos)
		lastpos.Y += v.Space.Height
		lastpos.Y += shift
	}
}

func AddOneNextToAnother(position engi.Point, shift float32, entities ...*Entity) {
	lastpos := position
	for _, v := range entities {
		v.PosSet(lastpos)
		lastpos.X += v.Space.Width
		lastpos.X += shift
	}
}
