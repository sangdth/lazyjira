package main

import (
	"log"

	ui "github.com/awesome-gocui/gocui"
	config "github.com/gookit/config/v2"
)

// RelativeSize returns the relative size of the terminal window view
// the first int is 30% of the width, the second is the 70% of the height
func relativeSize(g *ui.Gui) (int, int) {
	tw, th := g.Size()

	return (tw * 3) / 10, (th * 7) / 10
}

// TODO: Make helpers that do all the calculation like
// center vertically and horizontally (similar to margin auto in css)
// maybe something display flex could be great he he
func layout(g *ui.Gui) error {
	tw, th := g.Size()
	rw, rh := relativeSize(g)

	if !config.Exists(ServerKey) {
		if err := createPromptView(g, InsertServerTitle); err != nil {
			log.Panicln("Error while inserting server", err)
		}
	}

	if !config.Exists(UsernameKey) {
		if err := createPromptView(g, InsertUsernameTitle); err != nil {
			log.Panicln("Error while inserting username", err)
		}
	}

	if _, err := g.SetView(ProjectsView, 0, 0, rw, th-rh, 0); err != nil {
		log.Panicln("Cannot update view", err)
	}

	if _, err := g.View(StatusesView); err == nil {
		_, err := g.SetView(StatusesView, 0, 0, rw, th-rh, 0)
		if err != nil && err != ui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.View(PromptView); err == nil {
		_, err := g.SetView(PromptView, tw/6, (th/2)-8, (tw*5)/6, (th/2)-6, 0)
		if err != nil && err != ui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.View(AlertView); err == nil {
		_, err := g.SetView(AlertView, tw/6, (th/2)-10, (tw*5)/6, (th/2)-4, 0)
		if err != nil && err != ui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView(IssuesView, 0, th-rh+1, rw, th-3, 0); err != nil {
		log.Panicln("Cannot update view", err)
	}

	if _, err := g.SetView(DetailsView, rw+1, 0, tw-1, th-3, 0); err != nil {
		log.Panicln("Cannot update view", err)
	}

	return nil
}
