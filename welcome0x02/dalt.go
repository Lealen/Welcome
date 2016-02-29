package main

import (
	"image/color"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Dalt struct{}

func (c *Dalt) Preload() {
}

func (c *Dalt) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&engi.MouseSystem{})

	w.AddSystem(&DaltSystem{})

	Dalt := ecs.NewEntity([]string{"RenderSystem", "DaltSystem"})
	w.AddEntity(Dalt)
}

func (*Dalt) Hide()        {}
func (*Dalt) Show()        {}
func (*Dalt) Type() string { return "Dalt" }

type DaltSystem struct {
	*ecs.System

	mediumFont *engi.Font

	world *ecs.World

	red,
	green,
	blue uint8

	redc,
	greenc,
	bluec int8

	dragging *Entity
	lastpos  engi.Point

	changedt float32
}

func (DaltSystem) Type() string {
	return "DaltSystem"
}

func (c *DaltSystem) New(w *ecs.World) {
	c.System = ecs.NewSystem()
	c.world = w

	menufontname := "04b.ttf"

	c.mediumFont = GetFont(menufontname, 24, color.Transparent, color.White)

	NewEntity("Dalthudbg", []string{"RenderSystem", "MouseSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}),
		Position: engi.Point{X: 0, Y: 0},
		Width:    engi.Width(),
		Height:   engi.Height(),
		Priority: engi.HUDGround,
		OnWindowResize: func(e *Entity) {
			e.Space.Width = engi.Width()
			e.Space.Height = engi.Height()
			e.Render.SetScale(engi.Point{X: engi.Width(), Y: engi.Height()})
		},
	})

	var size float32 = 300

	NewEntity("Normal Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: c.red,
			G: c.green,
			B: c.blue,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 - size - 2, Y: engi.Height()/2 - size/2},
		Width:    size,
		Height:   size,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: c.red,
				G: c.green,
				B: c.blue,
				A: 255,
			}))
		},
	})

	var daltred, daltgreen, daltblue uint8

	NewEntity("Color Blind Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: c.red,
			G: c.green,
			B: c.blue,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 + 2, Y: engi.Height()/2 - size/2},
		Width:    size,
		Height:   size,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			var nasilenie float32 = 0.66 //0.5 //0.75

			daltred, daltgreen, daltblue = c.red, c.green, c.blue

			daltred = uint8(float32(c.red) /* * 0.5 */)
			daltgreen = uint8(float32(c.green) * nasilenie)
			daltblue = uint8(float32(c.blue) * nasilenie)

			proc := float32(c.red) / 255.0
			daltgreen += uint8(256 * (1 - nasilenie) * proc)
			daltblue += uint8(256 * (1 - nasilenie) * proc)

			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: daltred,
				G: daltgreen,
				B: daltblue,
				A: 255,
			}))
		},
	})

	NewEntity("Red Color Blind Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: c.red,
			G: 0,
			B: 0,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 + size + 4, Y: engi.Height()/2 - size/2},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: daltred,
				G: 0,
				B: 0,
				A: 255,
			}))
		},
	})

	NewEntity("Green Color Blind Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: c.green,
			B: 0,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 + size + 4, Y: engi.Height()/2 - size/6},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: 0,
				G: daltgreen,
				B: 0,
				A: 255,
			}))
		},
	})

	NewEntity("Blue Color Blind Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: 0,
			B: c.blue,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 + size + 4, Y: engi.Height()/2 + size/6},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: 0,
				G: 0,
				B: daltblue,
				A: 255,
			}))
		},
	})

	NewEntity("Red Color Normal Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: c.red,
			G: 0,
			B: 0,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 - size - size/3 - 4, Y: engi.Height()/2 - size/2},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: c.red,
				G: 0,
				B: 0,
				A: 255,
			}))
		},
	})

	NewEntity("Green Color Normal Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: c.green,
			B: 0,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 - size - size/3 - 4, Y: engi.Height()/2 - size/6},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: 0,
				G: c.green,
				B: 0,
				A: 255,
			}))
		},
	})

	NewEntity("Blue Color Normal Vision", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture: GetTextureColor(color.RGBA{
			R: 0,
			G: 0,
			B: c.blue,
			A: 255,
		}),
		Position: engi.Point{X: engi.Width()/2 - size - size/3 - 4, Y: engi.Height()/2 + size/6},
		Width:    size / 3,
		Height:   size / 3,
		Priority: engi.HUDGround + 1,
		OnUpdate: func(e *Entity, dt float32) {
			e.Render.SetDrawable(GetTextureColor(color.RGBA{
				R: 0,
				G: 0,
				B: c.blue,
				A: 255,
			}))
		},
	})

	NewEntity("testtext", []string{"RenderSystem"}, c.world, &EntityDefaults{
		Texture:  c.mediumFont.Render("text string"),
		Position: engi.Point{X: 40, Y: 40},
		Priority: engi.HUDGround + 1,
	})

}

func (c *DaltSystem) Update(entity *ecs.Entity, dt float32) {
	UpdateEntities(dt)

	c.changedt += dt

	if c.changedt > 1/64 {
		c.changedt -= 1 / 64
		switch {
		case c.red == 0 && c.green == 0 && c.blue == 0:
			c.redc = 0
			c.greenc = 0
			c.bluec = 1
		case c.red == 0 && c.green == 0 && c.blue == 255:
			c.redc = 0
			c.greenc = 1
			c.bluec = 0
		case c.red == 0 && c.green == 255 && c.blue == 255:
			c.redc = 0
			c.greenc = 0
			c.bluec = -1
		case c.red == 0 && c.green == 255 && c.blue == 0:
			c.redc = 1
			c.greenc = 0
			c.bluec = 0
		case c.red == 255 && c.green == 255 && c.blue == 0:
			c.redc = 0
			c.greenc = 0
			c.bluec = 1
		case c.red == 255 && c.green == 255 && c.blue == 255:
			c.redc = 0
			c.greenc = -1
			c.bluec = 0
		case c.red == 255 && c.green == 0 && c.blue == 255:
			c.redc = 0
			c.greenc = 0
			c.bluec = -1
		case c.red == 255 && c.green == 0 && c.blue == 0:
			c.redc = -1
			c.greenc = 0
			c.bluec = 0
		}

		if c.redc == 1 {
			c.red++
		} else if c.redc == -1 {
			c.red--
		}
		if c.greenc == 1 {
			c.green++
		} else if c.greenc == -1 {
			c.green--
		}
		if c.bluec == 1 {
			c.blue++
		} else if c.bluec == -1 {
			c.blue--
		}
	}
}
