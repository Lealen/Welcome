package main

import (
	"fmt"

	"github.com/Lealen/engi"
	"github.com/Lealen/engi/ecs"
)

type Game struct{}

func (c *Game) Preload() {}

func (c *Game) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&PlayerSystem{})
	w.AddSystem(&DebugSystem{})

	player := ecs.NewEntity([]string{"RenderSystem", "PlayerSystem", "DebugSystem"})
	fmt.Println("czym jest player: %#v", player)
	w.AddEntity(player)

	CameraSetX(engi.Width() / 2)
	CameraSetY(engi.Height() / 2)
	CameraSetZ(1.0)

	//fmt.Printf("Game\t world: %p:\n%#v\n\n", w, w)
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }
