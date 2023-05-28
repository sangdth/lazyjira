package main

import (
	"fmt"
	"log"

	ui "github.com/awesome-gocui/gocui"
)

type Category string

const (
	PROMPT Category = "prompt"
	ALERT  Category = "alert"
)

// Dialog is just like Dialog, instead it's mostly used for displaying an input
// or a message to the user.
type Dialog struct {
	*ui.View
	message  string
	value    string
	category Category
}

type CreateDialogOptions struct {
	title   string
	content string
	value   string
}

// CreateDialog initializes a Dialog object with an existing View by applying some
// basic configuration
func CreateDialog(v *ui.View, c Category) *Dialog {
	dialog := &Dialog{}
	dialog.View = v
	dialog.category = c

	v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝'}

	switch c {
	case PROMPT:
		v.Editable = true
		v.FrameColor = ui.ColorBlue
		v.TitleColor = ui.ColorBlue

	case ALERT:
		v.FrameColor = ui.ColorRed
		v.TitleColor = ui.ColorRed

	default:
		v.FrameColor = ui.ColorGreen
		v.TitleColor = ui.ColorGreen
	}

	return dialog
}

// IsEmpty indicates whether a dialog has items or not
func (d *Dialog) IsEmpty() bool {
	return d.Length() == 0
}

// Focus hightlights the View of the current Dialog
func (d *Dialog) Focus(g *ui.Gui) {

	switch d.category {
	case PROMPT:
		g.Cursor = true
		d.Editable = true
		d.FrameColor = ui.ColorBlue
		d.TitleColor = ui.ColorBlue

	case ALERT:
		g.Cursor = false
		d.FrameColor = ui.ColorRed
		d.TitleColor = ui.ColorRed

	default:
		d.FrameColor = ui.ColorDefault
		d.TitleColor = ui.ColorDefault
	}

	_, err := g.SetCurrentView(d.Name())
	if err != nil {
		log.Panicln("Error on Focus", err)
	}
}

// Unfocus is used to remove highlighting from the current dialog
func (d *Dialog) Unfocus() {
	d.FrameColor = ui.ColorDefault
	d.TitleColor = ui.ColorDefault
}

/**
 * Set the title of the View and display paging information of the
 * dialog if there are more than one pages
 */
func (d *Dialog) SetTitles(title string, subtitle string) {
	d.Title = title
	d.Subtitle = subtitle
}

/**
 * Sometimes the value need to be passed to next action through the Alert confirmation
 */
func (d *Dialog) SetValue(v string) {
	d.value = v
}

// SetItems will (re)evaluates the dialog's items with the given data and redraws
// the View
func (d *Dialog) SetContent(content string) {
	if _, err := fmt.Fprintln(d.View, content); err != nil {
		log.Panicln("Error on SetContent", err)
	}
}

// ResetCursor puts the cirson back at the beginning of the View
func (d *Dialog) ResetCursor() {
	err := d.SetCursor(0, 0)
	if err != nil {
		log.Panicln("Error in ResetCursor", err)
	}
}

// currentCursorY returns the current Y of the cursor
func (d *Dialog) CurrentCursorY() int {
	_, y := d.Cursor()

	return y
}

// height ewturns the current height of the View
func (d *Dialog) Height() int {
	_, y := d.Size()

	return y - 1
}

// width ewturns the current width of the View
func (d *Dialog) Width() int {
	x, _ := d.Size()

	return x - 1
}

// length returns the length of the dialog
func (d *Dialog) Length() int {
	return len(d.message)
}
