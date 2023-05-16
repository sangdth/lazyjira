package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	issues, _ := ListIssuesByProjectCode("FF")

	// Initialize the gocui library
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	g.Cursor = true

	// Set up the main screen and keybindings
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	issuesView, _ := g.SetView("issues", 0, 0, 50, len(issues)+2)
	issuesView.Title = " Issues "
	issuesView.Highlight = true
	issuesView.SelBgColor = gocui.ColorGreen
	issuesView.SelFgColor = gocui.ColorBlack

	for _, issue := range issues {
		// Extract relevant information from the issue
		key := issue.Key
		summary := issue.Fields.Summary

		// Format the row
		row := fmt.Sprintf("%-10s %s", key, summary)

		// Add the row to the issues view
		fmt.Fprintln(issuesView, row)
	}

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Set keybindings for the views
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("issues", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	if err := g.SetKeybinding("projects", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("projects", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("issues", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("issues", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("issues", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("issues", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
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

// The layout function sets up the main screen layout
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("projects", 0, 0, 50, 20); err != nil {
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

	if v, err := g.SetView("issues", 0, 21, 50, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Issues "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView("details", 51, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Details "
		v.Wrap = true
		if _, err := g.SetCurrentView("projects"); err != nil {
			return err
		}
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "projects" {
		_, err := g.SetCurrentView("issues")
		return err
	}
	_, err := g.SetCurrentView("projects")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
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

func cursorUp(g *gocui.Gui, v *gocui.View) error {
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
