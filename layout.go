package main

import (
	"log"

	ui "github.com/awesome-gocui/gocui"
)

// RelativeSize returns the relative size of the terminal window view
// the first int is 30% of the width, the second is the 70% of the height
func relativeSize(g *ui.Gui) (int, int) {
	tw, th := g.Size()

	return (tw * 3) / 10, (th * 7) / 10
}

func layout(g *ui.Gui) error {
	tw, th := g.Size()
	rw, rh := relativeSize(g)

	_, err := g.SetView(ProjectsView, 0, 0, rw, th-rh, 0)
	if err != nil {
		log.Panicln("Cannot update view", err)
	}

	if _, err := g.View(StatusesView); err == nil {
		_, err = g.SetView(StatusesView, 0, 0, rw, th-rh, 0)
		if err != nil && err != ui.ErrUnknownView {
			return err
		}
	}

	if _, err = g.View(PromptView); err == nil {
		_, err = g.SetView(PromptView, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1, 0)
		if err != nil && err != ui.ErrUnknownView {
			return err
		}
	}

	_, err = g.SetView(IssuesView, 0, th-rh+1, rw, th-3, 0)
	if err != nil {
		log.Panicln("Cannot update view", err)
	}

	_, err = g.SetView(DetailsView, rw+1, 0, tw-1, th-3, 0)
	if err != nil {
		log.Panicln("Cannot update view", err)
	}

	return nil
}
