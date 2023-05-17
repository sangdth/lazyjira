package main

import (
	"github.com/jroimartin/gocui"
)

const (
	allViewsKey     = ""
	projectsViewKey = "projects"
	issuesViewKey   = "issues"
	detailsViewKey  = "details"
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(issuesViewKey, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(projectsViewKey, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	// PROJECTS VIEW
	// Use j/k and arrow down/up to navigate in projects view
	if err := g.SetKeybinding(projectsViewKey, 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(projectsViewKey, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(projectsViewKey, 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(projectsViewKey, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	// ISSUES VIEW
	// Use j/k and arrow down/up to navigate in issues view
	if err := g.SetKeybinding(issuesViewKey, 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(issuesViewKey, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(issuesViewKey, 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(issuesViewKey, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	// ALL VIEWS
	// Use Ctrl-c or q to quit
	if err := g.SetKeybinding(allViewsKey, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(allViewsKey, 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	// if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
	// 	return err
	// }

	// if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
	// 	return err
	// }
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "projects" {
		_, err := g.SetCurrentView("issues")
		return err
	}
	_, err := g.SetCurrentView("projects")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
