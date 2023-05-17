package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	sideBarWidth      = 50
	projectViewHeight = 20
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(projectsViewKey, 0, 0, sideBarWidth, projectViewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Projects "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		savedProjects := GetSavedProjects()
		for _, project := range savedProjects {
			fmt.Fprintln(v, project)
		}
	}

	if v, err := g.SetView(issuesViewKey, 0, projectViewHeight+1, sideBarWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Issues "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView(detailsViewKey, sideBarWidth+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Details "
		v.Wrap = true
		if _, err := g.SetCurrentView(projectsViewKey); err != nil {
			return err
		}
	}

	return nil
}
