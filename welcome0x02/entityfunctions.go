package main

import "github.com/Lealen/engi"

var entititesFunctions []*Entity

var lastwindowwidth, lastwindowheight float32
var PreviousMousePosX, PreviousMousePosY float32

func UpdateEntities() {
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
		if v.Mouse.Clicked {
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
		if v.OnDragged != nil && v.Mouse.Dragged {
			v.OnDragged(v)
		}
		//		if v.OnRightClicked != nil && v.Mouse.Rightclicked {
		//			v.OnRightClicked(v)
		//		}
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
