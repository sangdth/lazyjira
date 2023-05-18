package main

import (
	ui "github.com/awesome-gocui/gocui"
)

func createStatusView(g *ui.Gui) error {
	_, th := g.Size()
	rw, rh := relativeSize(g)

	v, err := g.SetView(StatusesView, 0, 0, rw, th-rh, 0)
	if err != nil && err != ui.ErrUnknownView {
		return err
	}
	StatusesList = CreateList(v, false)
	StatusesList.Title = MakeProjectTabNames(StatusesView)

	_, err = g.SetCurrentView(StatusesView)

	return err
}
