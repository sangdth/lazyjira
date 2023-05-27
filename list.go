package main

import (
	"fmt"
	"log"

	ui "github.com/awesome-gocui/gocui"
)

// Page is used to hold info about a list based view
type Page struct {
	offset, limit int
}

// List overlads the gocui.View by implementing list specific functionalitys
type List struct {
	*ui.View
	code      string
	title     string
	items     []string
	pages     []Page
	pageIndex int
	ordered   bool
}

// CreateList initializes a List object with an existing View by applying some
// basic configuration
func CreateList(v *ui.View, ordered bool) *List {
	list := &List{}
	list.View = v
	list.Autoscroll = true
	list.ordered = ordered

	return list
}

// IsEmpty indicates whether a list has items or not
func (l *List) IsEmpty() bool {
	return l.length() == 0
}

// Focus hightlights the View of the current List
func (l *List) Focus(g *ui.Gui) {
	l.Highlight = true
	l.SelFgColor = ui.ColorBlack
	l.SelBgColor = ui.ColorGreen
	l.FrameColor = ui.ColorGreen
	l.TitleColor = ui.ColorGreen
	_, err := g.SetCurrentView(l.Name())
	if err != nil {
		log.Panicln("Error on SetCurrentView", err)
	}
}

// Unfocus is used to remove highlighting from the current list
func (l *List) Unfocus() {
	l.FrameColor = ui.ColorDefault
	l.TitleColor = ui.ColorDefault
	l.Highlight = false
}

// Reset zeros the list's slices out and clears the underlying View
func (l *List) Reset() {
	l.items = make([]string, 0)
	l.pages = []Page{}
	l.Clear()
	l.ResetCursor()
}

// Change the project code means old data will be gone
func (l *List) SetCode(code string) {
	// only do if code is new
	if l.code != code {
		l.code = code
		l.Reset()
	}
}

// SetTitle will set the title of the View and display paging information of the
// list if there are more than one pages
func (l *List) SetTitle(title string) {
	l.title = title

	if l.pagesNum() > 1 {
		l.Title = fmt.Sprintf("%d/%d - %s", l.currPageNum(), l.pagesNum(), title)
	} else {
		l.Title = title
	}
}

// SetItems will (re)evaluates the list's items with the given data and redraws
// the View
func (l *List) SetItems(data []string) {
	l.items = data
	l.ResetPages()
	err := l.Draw()
	if err != nil {
		log.Panicln("Error on SetItems", err)
	}
}

// AddItem appends a given item to the existing list
func (l *List) AddItem(g *ui.Gui, item string) {
	l.items = append(l.items, item)
	l.ResetPages()
	if err := l.Draw(); err != nil {
		log.Panicln("Error on AddItem", err)
	}
}

func (l *List) UpdateCurrentItem(item string) {
	page := l.currPage()
	data := l.items[page.offset : page.offset+page.limit]

	data[l.currentCursorY()] = item
}

// Draw calculates the pages and draws the first one
func (l *List) Draw() error {
	if l.IsEmpty() {
		return nil
	}
	return l.displayPage(0)
}

// Draw calculates the pages and draws the first one
func (l *List) DrawCurrentPage() error {
	if l.IsEmpty() {
		return nil
	}
	return l.displayPage(l.pageIndex)
}

// MoveDown moves the cursor to the line below or the next page if any
func (l *List) MoveDown() error {
	if l.IsEmpty() {
		return nil
	}
	y := l.currentCursorY() + 1
	if l.atBottomOfPage() {
		y = 0
		if l.hasMultiplePages() {
			return l.displayPage(l.nextPageIdx())
		}
	}
	err := l.SetCursor(0, y)
	if err != nil {
		return err
	}

	return nil
}

// MoveUp moves the cursor to the line above on the previous page if any
func (l *List) MoveUp() error {
	if l.IsEmpty() {
		return nil
	}
	y := l.currentCursorY() - 1
	if l.atTopOfPage() {
		y = l.pages[l.prevPageIdx()].limit - 1
		if l.hasMultiplePages() {
			return l.displayPage(l.prevPageIdx())
		}
	}

	err := l.SetCursor(0, y)
	if err != nil {
		return err
	}

	return nil
}

// MovePgDown displays the next page
func (l *List) MovePgDown() error {
	if l.IsEmpty() {
		return nil
	}

	err := l.displayPage(l.nextPageIdx())
	if err != nil {
		return err
	}

	return l.SetCursor(0, 0)
}

// MovePgUp displays the previous page
func (l *List) MovePgUp() error {
	if l.IsEmpty() {
		return nil
	}
	err := l.displayPage(l.prevPageIdx())
	if err != nil {
		log.Panicln(err)
	}

	return l.SetCursor(0, 0)
}

// CurrentItem returns the currently selected item of the list no matter what
// page is being displayed
func (l *List) CurrentItem() string {
	if l.IsEmpty() {
		return ""
	}
	page := l.currPage()
	data := l.items[page.offset : page.offset+page.limit]

	return data[l.currentCursorY()]
}

// ResetCursor puts the cirson back at the beginning of the View
func (l *List) ResetCursor() {
	err := l.SetCursor(0, 0)
	if err != nil {
		log.Panicln("Error in ResetCursor", err)
	}
}

// ResetPages (re)calculates the pages data based on the current length of the
// list and the current height of the View
func (l *List) ResetPages() {
	l.pages = []Page{}
	for offset := 0; offset < l.length(); offset += l.height() {
		limit := l.height()
		if offset+limit > l.length() {
			limit = l.length() % l.height()
		}
		l.pages = append(l.pages, Page{offset, limit})
	}
}

// currPageNum returns the current page number to display
func (l *List) currPageNum() int {
	if l.IsEmpty() {
		return 0
	}
	return l.pageIndex + 1
}

// currentCursorY returns the current Y of the cursor
func (l *List) currentCursorY() int {
	_, y := l.Cursor()

	return y
}

// currPage returns the current page being displayd
func (l *List) currPage() Page {
	return l.pages[l.pageIndex]
}

// height ewturns the current height of the View
func (l *List) height() int {
	_, y := l.Size()

	return y - 1
}

// width ewturns the current width of the View
func (l *List) width() int {
	x, _ := l.Size()

	return x - 1
}

// length returns the length of the list
func (l *List) length() int {
	return len(l.items)
}

// pageNum returns the number of the pages
func (l *List) pagesNum() int {
	return len(l.pages)
}

// nextPageIdx returns the index of the next page to be displayed circularlt
func (l *List) nextPageIdx() int {
	return (l.pageIndex + 1) % l.pagesNum()
}

// prevPageIdx returns the index of the prev page to be displayed circularlt
func (l *List) prevPageIdx() int {
	pidx := (l.pageIndex - 1) % l.pagesNum()
	if l.pageIndex == 0 {
		pidx = l.pagesNum() - 1
	}
	return pidx
}

// sidplayItem displays the text of the item with index i and fills with spaces
// the remaining space until the border of the View
func (l *List) displayItem(i int) string {
	item := fmt.Sprint(l.items[i])
	sp := spaces(l.width() - len(item) + 1)
	if l.ordered {
		return fmt.Sprintf("%2d. %v%s", i+1, item, sp)
	} else {
		return fmt.Sprintf("%s%s", item, sp)
	}
}

// displayPage resets the currentPageIdx and displays the requested page
func (l *List) displayPage(p int) error {
	l.Clear()
	l.pageIndex = p
	page := l.pages[l.pageIndex]
	for i := page.offset; i < page.offset+page.limit; i++ {
		if _, err := fmt.Fprintln(l.View, l.displayItem(i)); err != nil {
			return err
		}
	}
	l.SetTitle(l.title)

	return nil
}

// atBottomOfPage determines whether the cursor is at the bottom of the current page
func (l *List) atBottomOfPage() bool {
	return l.currentCursorY() == l.currPage().limit-1
}

// atTopOfPage determines whether the cursor is at the top of the current page
func (l *List) atTopOfPage() bool {
	return l.currentCursorY() == 0
}

// hasMultiplePages determines whether there is more than one page to be displayed
func (l *List) hasMultiplePages() bool {
	return l.pagesNum() > 1
}
