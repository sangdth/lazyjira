package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func main() {
	// Initialize the gocui library
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("Error initializing gocui library: ", err)
		return
	}
	defer g.Close()

	// Set up the main screen and keybindings
	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		fmt.Println("Error setting up keybindings: ", err)
		return
	}

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println("Error in main event loop: ", err)
	}
}

// The layout function sets up the main screen layout
func layout(g *gocui.Gui) error {
	// TODO: Implement the layout function
	return nil
}

// The keybindings function sets up the main keybindings
func keybindings(g *gocui.Gui) error {
	// TODO: Implement the keybindings function
	return nil
}
