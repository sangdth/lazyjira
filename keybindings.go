package main

import (
	"log"

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
	if err := g.SetKeybinding(ProjectsView, 'j', ui.ModNone, ListDown); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowDown, ui.ModNone, ListDown); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(ProjectsView, 'k', ui.ModNone, ListUp); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeyArrowUp, ui.ModNone, ListUp); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(ProjectsView, ui.KeySpace, ui.ModNone, OnSelectProject); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	// ISSUES VIEW
	if err := g.SetKeybinding(IssuesView, 'j', ui.ModNone, ListDown); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowDown, ui.ModNone, ListDown); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(IssuesView, 'k', ui.ModNone, ListUp); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(IssuesView, ui.KeyArrowUp, ui.ModNone, ListUp); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	// ALL VIEWS
	if err := g.SetKeybinding(AllViews, ui.KeyCtrlC, ui.ModNone, Quit); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(AllViews, 'q', ui.ModNone, Quit); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(AllViews, ui.KeyEnter, ui.ModNone, OnEnter); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	return nil
}
