package main

import (
	ui "github.com/jroimartin/gocui"
)

const (
	AllViews     = ""
	ProjectsView = "projects"
	IssuesView   = "issues"
	DetailsView  = "details"
)

func keybindings(g *ui.Gui) error {
	if err := g.SetKeybinding(IssuesView, ui.KeyTab, ui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyTab, ui.ModNone, nextView); err != nil {
		return err
	}

	// PROJECTS VIEW
	// Use j/k and arrow down/up to navigate in projects view
	if err := g.SetKeybinding(ProjectsView, 'j', ui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowDown, ui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, 'k', ui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowUp, ui.ModNone, cursorUp); err != nil {
		return err
	}

	// ISSUES VIEW
	// Use j/k and arrow down/up to navigate in issues view
	if err := g.SetKeybinding(IssuesView, 'j', ui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowDown, ui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, 'k', ui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowUp, ui.ModNone, cursorUp); err != nil {
		return err
	}

	// ALL VIEWS
	// Use Ctrl-c or q to quit
	if err := g.SetKeybinding(AllViews, ui.KeyCtrlC, ui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(AllViews, 'q', ui.ModNone, quit); err != nil {
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

func nextView(g *ui.Gui, v *ui.View) error {
	if v == nil || v.Name() == "projects" {
		_, err := g.SetCurrentView("issues")
		return err
	}
	_, err := g.SetCurrentView("projects")
	return err
}

func cursorDown(g *ui.Gui, v *ui.View) error {
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

func cursorUp(g *ui.Gui, v *ui.View) error {
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

func quit(g *ui.Gui, v *ui.View) error {
	return ui.ErrQuit
}
