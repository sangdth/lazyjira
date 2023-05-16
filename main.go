package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	// Initialize the gocui library
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	// Set up the main screen and keybindings
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Start the main event loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// The layout function sets up the main screen layout
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("colors", maxX/2-7, maxY/2-12, maxX/2+7, maxY/2+13); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for i := 0; i <= 7; i++ {
			for _, j := range []int{1, 4, 7} {
				fmt.Fprintf(v, "lazy he he \033[3%d;%dmcolors!\033[0m\n", i, j)
			}
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
