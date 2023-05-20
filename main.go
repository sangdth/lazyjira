package main

import (
	"log"

	ui "github.com/awesome-gocui/gocui"
)

const (
	AllViews     = ""
	ProjectsView = "projects"
	StatusesView = "statuses"
	IssuesView   = "issues"
	DetailsView  = "details"
	PromptView   = "prompt"
)

var (
	ProjectsList *List
	StatusesList *List
	IssuesList   *List
	Details      *ui.View
)

func main() {
	InitConfig()

	// Initialize the gocui library
	g, err := ui.NewGui(ui.OutputNormal, true)
	if err != nil {
		log.Panicln("Failed to initialize GUI", err)
	}

	defer g.Close()

	g.Cursor = false

	// Set up the main screen and keybindings
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln("Failed to attach keybindings", err)
	}

	tw, th := g.Size()
	rw, rh := relativeSize(g)

	v, err := g.SetView(ProjectsView, 0, 0, rw, th-rh, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create view", err)
	}
	ProjectsList = CreateList(v, false)
	ProjectsList.Title = makeTabNames(ProjectsView)
	ProjectsList.Focus(g)

	g.Update(func(g *ui.Gui) error {
		LoadProjects()

		return nil
	})

	v, err = g.SetView(IssuesView, 0, th-rh+1, rw, th-3, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create view", err)
	}
	IssuesList = CreateList(v, false)
	IssuesList.Title = " Issues "

	Details, err = g.SetView(DetailsView, rw+1, 0, tw-1, th-3, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Failed to create Details view", err)
	}
	Details.Title = " Details "
	Details.Wrap = true

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != ui.ErrQuit {
		log.Panicln(err)
	}
}
