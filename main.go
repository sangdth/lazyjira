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

	issuesV, _ := g.SetView(issuesViewKey, 0, 0, 50, len(issues)+2)
	issuesV.Title = " Issues "
	issuesV.Highlight = true
	issuesV.SelBgColor = gocui.ColorGreen
	issuesV.SelFgColor = gocui.ColorBlack

	for _, issue := range issues {
		// Extract relevant information from the issue
		key := issue.Key
		summary := issue.Fields.Summary

		// Format the row
		row := fmt.Sprintf("%-10s %s", key, summary)

		// Add the row to the issues view
		fmt.Fprintln(issuesV, row)
	}

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
