package main

import (
	"log"
	"strings"

	ui "github.com/awesome-gocui/gocui"
)

func CreateStatusView(g *ui.Gui) error {
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

func ListUp(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case ProjectsView:
		if err := ProjectsList.MoveUp(); err != nil {
			log.Println("Error on ProjectsList.MoveUp()", err)
			return err
		}
	case StatusesView:
		if err := StatusesList.MoveUp(); err != nil {
			log.Println("Error on StatusesList.MoveUp()", err)
			return err
		}
	case IssuesView:
		if err := IssuesList.MoveUp(); err != nil {
			log.Println("Error on IssuesList.MoveUp()", err)
			return err
		}
	}
	return nil
}

func ListDown(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case ProjectsView:
		if err := ProjectsList.MoveDown(); err != nil {
			log.Println("Error on SitesList.MoveDown()", err)
			return err
		}
	case StatusesView:
		if err := StatusesList.MoveDown(); err != nil {
			log.Println("Error on StatusesList.MoveDown()", err)
			return err
		}
	case IssuesView:
		if err := IssuesList.MoveDown(); err != nil {
			log.Println("Error on NewsList.MoveDown()", err)
			return err
		}
	}
	return nil
}

func ChangeView(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case ProjectsView:
		if v == ProjectsList.View {
			ProjectsList.Unfocus()
		}
		if strings.Contains(IssuesList.Title, "bookmarks") {
			g.SelFgColor = ui.ColorMagenta | ui.AttrBold
		}

		IssuesList.Focus(g)
	case IssuesView:
		ProjectsList.Focus(g)
		IssuesList.Unfocus()
	}

	return nil
}

func SwitchProjectTab(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case StatusesView:
		ProjectsList.Focus(g)
		StatusesList.Unfocus()
		g.DeleteView(StatusesView)

	case ProjectsView:
		if err := CreateStatusView(g); err == nil {
			OnEnter(g, v)
			ProjectsList.Unfocus()
			StatusesList.Focus(g)
		} else {
			log.Panicln("Error on createStatusView()", err)
		}
	}

	return nil
}

// Pressing Spacebar will trigger this one
func OnSelectProject(g *ui.Gui, v *ui.View) error {
	currentItem := ProjectsList.CurrentItem()
	if currentItem == nil {
		return nil
	}

	IssuesList.Clear()

	err := FetchIssues(g, currentItem.(string))

	return err
}

func OnEnter(g *ui.Gui, v *ui.View) error {
	currentItem := ProjectsList.CurrentItem()
	if currentItem == nil {
		return nil
	}

	projectCode := currentItem.(string)

	if IssuesList.IsEmpty() || IssuesList.code != projectCode {
		err := OnSelectProject(g, v)
		if err != nil {
			log.Println("Error on OnSelectProject", err)
		}
	}

	if err := CreateStatusView(g); err == nil {
		err := FetchStatuses(g, projectCode)
		if err != nil {
			log.Println("Error on FetchStatuses", err)
		}
	}

	return nil
}

func MakeProjectTabNames(name string) string {
	switch name {

	case ProjectsView:
		return " Projects "

	case StatusesView:
		return " Projects > Statuses "
	}

	return "Something went wrong in making name"
}

func Quit(g *ui.Gui, v *ui.View) error {
	return ui.ErrQuit
}
