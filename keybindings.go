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
	if err := g.SetKeybinding(AllViews, ui.KeyTab, ui.ModNone, SwitchView); err != nil {
		return err
	}

	// PROJECTS VIEW
	// Use j/k and arrow down/up to navigate in projects view
	if err := g.SetKeybinding(ProjectsView, 'j', ui.ModNone, ListDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowDown, ui.ModNone, ListDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, 'k', ui.ModNone, ListUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowUp, ui.ModNone, ListUp); err != nil {
		return err
	}

	// ISSUES VIEW
	// Use j/k and arrow down/up to navigate in issues view
	if err := g.SetKeybinding(IssuesView, 'j', ui.ModNone, ListDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowDown, ui.ModNone, ListDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, 'k', ui.ModNone, ListUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowUp, ui.ModNone, ListUp); err != nil {
		return err
	}

	// ALL VIEWS
	// Use Ctrl-c or q to quit
	if err := g.SetKeybinding(AllViews, ui.KeyCtrlC, ui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(AllViews, 'q', ui.ModNone, Quit); err != nil {
		return err
	}

	return nil
}
