package main

import "github.com/Lealen/engi"

var entititesFunctions []*Entity

var lastwindowwidth, lastwindowheight float32
var PreviousMousePosX, PreviousMousePosY float32

func IsMouseOn(e *Entity) bool {
	mx := engi.Mouse.X
	my := engi.Mouse.Y
	if mx > e.Space.Position.X && mx < (e.Space.Position.X+e.Space.Width) &&
		my > e.Space.Position.Y && my < (e.Space.Position.Y+e.Space.Height) {
		return true
	}
	return false
}

func UpdateEntities() {
	for _, v := range entititesFunctions {
		if v.OnUpdate != nil {
			v.OnUpdate(v)
		}
	}

	if lastwindowwidth == 0 && lastwindowheight == 0 {
		lastwindowwidth = engi.WindowWidth()
		lastwindowheight = engi.WindowHeight()
	} else if lastwindowwidth != engi.WindowWidth() || lastwindowheight != engi.WindowHeight() {
		for _, v := range entititesFunctions {
			if v.OnWindowResize != nil {
				v.OnWindowResize(v)
			}
		}
		lastwindowwidth = engi.WindowWidth()
		lastwindowheight = engi.WindowHeight()
	}

	for _, v := range entititesFunctions {
		if v.Mouse == nil {
			continue
		}
		ison := IsMouseOn(v)
		if ison && engi.Keys.Get(engi.MouseButtonLeft).Down() {
			if !v.IsClicked {
				if v.OnPress != nil {
					v.OnPress(v)
				}
				v.IsClicked = true
			}
			if v.OnClicked != nil {
				v.OnClicked(v)
			}
		} else if v.IsClicked {
			if v.OnRelease != nil {
				v.OnRelease(v)
			}
			v.IsClicked = false
		}
		if v.OnDragged != nil && v.Mouse.Dragged && engi.Keys.Get(engi.MouseButtonLeft).Down() {
			v.OnDragged(v)
		}
		if ison && engi.Keys.Get(engi.MouseButtonRight).Down() {
			if !v.IsRightClicked {
				if v.OnRightPress != nil {
					v.OnRightPress(v)
				}
				v.IsRightClicked = true
			}
			if v.OnRightClicked != nil {
				v.OnRightClicked(v)
			}
		} else if v.IsRightClicked {
			if v.OnRightRelease != nil {
				v.OnRightRelease(v)
			}
			v.IsRightClicked = false
		}
		if v.OnEnter != nil && v.Mouse.Enter {
			v.OnEnter(v)
		}
		if v.OnLeave != nil && v.Mouse.Leave {
			v.OnLeave(v)
		}
	}

	PreviousMousePosX = engi.Mouse.X
	PreviousMousePosY = engi.Mouse.Y
}
