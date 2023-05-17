package main

import (
	"log"

	ui "github.com/jroimartin/gocui"
)

var (
	ProjectsList *List
	IssuesList   *List
	Details      *ui.View
)

func main() {
	// issues, _ := ListIssuesByProjectCode("FF")

	// Initialize the gocui library
	g, err := ui.NewGui(ui.OutputNormal)
	if err != nil {
		log.Panicln("Failed to initialize GUI", err)
	}

	defer g.Close()

	g.Cursor = false
	g.Highlight = true

	// Set up the main screen and keybindings
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln("Failed to attach keybindings", err)
	}

	// savedProjects := GetSavedProjects()
	// for _, project := range savedProjects {
	// 	fmt.Fprintln(v, project)
	// }

	tw, th := g.Size()
	rw, rh := relativeSize(g)

	v, err := g.SetView(ProjectsView, 0, 0, rw, th-rh)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create Projects view", err)
	}
	ProjectsList = CreateList(v, true)
	ProjectsList.Title = " Projects "
	ProjectsList.Focus(g)

	g.Update(func(g *ui.Gui) error {
		if err := LoadSites(); err != nil {
			log.Panicln("Error while loading projects", err)
		}
		log.Print("Loaded initial projects")
		return nil
	})

	v, err = g.SetView(IssuesView, 0, th-rh+1, rw, th-3)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create Issues view", err)
	}
	IssuesList = CreateList(v, true)
	IssuesList.Title = " Issues "

	Details, err = g.SetView(DetailsView, rw+1, 0, tw-1, th-3)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create Details view", err)
	}
	Details.Title = " Details "
	Details.Wrap = true

	// for _, issue := range issues {
	// 	// Extract relevant information from the issue
	// 	key := issue.Key
	// 	summary := issue.Fields.Summary

	// 	// Format the row
	// 	row := fmt.Sprintf("%-1s %s%*s", key, summary, 51-(len(key)+len(summary)), "")

	// 	// Add the row to the issues view
	// 	fmt.Fprintln(v, row)
	// }

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != ui.ErrQuit {
		log.Panicln(err)
	}
}
