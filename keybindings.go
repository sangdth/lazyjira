package main

import (
	"log"

	ui "github.com/awesome-gocui/gocui"
)

func keybindings(g *ui.Gui) error {
	if err := g.SetKeybinding(AllViews, ui.KeyTab, ui.ModNone, ChangeView); err != nil {
		return err
	}

	// PROJECTS VIEW
	if err := g.SetKeybinding(ProjectsView, 'a', ui.ModNone, AddProject); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(ProjectsView, 'd', ui.ModNone, RemoveProject); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
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
	if err := g.SetKeybinding(ProjectsView, ui.KeyEnter, ui.ModNone, SwitchProjectTab); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	// STATUSES VIEW
	if err := g.SetKeybinding(StatusesView, 'b', ui.ModNone, SwitchProjectTab); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(StatusesView, ui.KeyEsc, ui.ModNone, SwitchProjectTab); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(StatusesView, 'j', ui.ModNone, ListDown); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(StatusesView, 'k', ui.ModNone, ListUp); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(StatusesView, ui.KeySpace, ui.ModNone, ToggleStatus); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	// PROMPT VIEW
	if err := g.SetKeybinding(PromptView, ui.KeyEsc, ui.ModNone, CancelDialog); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(PromptView, ui.KeyEnter, ui.ModNone, SubmitPrompt); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}

	// ALERT VIEW
	if err := g.SetKeybinding(AlertView, ui.KeyEsc, ui.ModNone, CancelDialog); err != nil {
		log.Fatal("Failed to set keybindings", err)
	}
	if err := g.SetKeybinding(AlertView, ui.KeyEnter, ui.ModNone, SubmitAlert); err != nil {
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

	return nil
}
