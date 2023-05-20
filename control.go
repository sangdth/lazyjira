package main

import (
	"fmt"
	"log"
	"strings"

	ui "github.com/awesome-gocui/gocui"
	viper "github.com/spf13/viper"
)

func CreateStatusView(g *ui.Gui) error {
	_, th := g.Size()
	rw, rh := relativeSize(g)

	v, err := g.SetView(StatusesView, 0, 0, rw, th-rh, 0)
	if err != nil && err != ui.ErrUnknownView {
		return err
	}
	StatusesList = CreateList(v, false)
	StatusesList.Title = makeTabNames(StatusesView)

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
		if err := g.DeleteView(StatusesView); err != nil {
			return err
		}

	case ProjectsView:
		if err := CreateStatusView(g); err == nil {
			if err := OnEnter(g, v); err != nil {
				return err
			}
			ProjectsList.Unfocus()
			StatusesList.Focus(g)
		} else {
			log.Panicln("Error on createStatusView()", err)
		}
	}

	return nil
}

func ToggleStatus(g *ui.Gui, v *ui.View) error {
	currentItem := StatusesList.CurrentItem()
	if currentItem == nil {
		return nil
	}

	projectCode := IssuesList.code

	value := currentItem.(string)

	path := fmt.Sprintf("savedprojects.%s.statuses.%s", projectCode, value)

	if viper.IsSet(path) {
		currentValue := viper.GetBool(path)
		viper.Set(path, !currentValue)
	} else {
		viper.Set(path, true)
	}

	if err := viper.WriteConfig(); err != nil {
		return err
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

	IssuesList.Title = " Issues | Fetching... "

	// Can not nest the update
	g.Update(func(g *ui.Gui) error {
		if err := FetchIssues(g, currentItem.(string)); err != nil {
			return err
		}

		IssuesList.Title = " Issues "

		return nil
	})

	return nil
}

// When pressing Enter, the Issues list might be empty, so we need to fetch it again
func OnEnter(g *ui.Gui, v *ui.View) error {
	currentItem := ProjectsList.CurrentItem()
	if currentItem == nil {
		return nil
	}

	projectCode := currentItem.(string)

	if err := CreateStatusView(g); err != nil {
		return err
	}

	IssuesList.Title = " Issues | Fetching... "
	StatusesList.Title = " Projects > Statuses | Fetching... "

	// Can not nest the update
	g.Update(func(g *ui.Gui) error {
		if IssuesList.IsEmpty() || IssuesList.code != projectCode {
			if err := FetchIssues(g, projectCode); err != nil {
				return err
			}
		}

		oldStatuses := GetSavedStatusesByProjectCode(projectCode)
		if len(oldStatuses) > 0 {
			StatusesList.SetItems(oldStatuses)
		} else {
			if err := FetchStatuses(g, projectCode); err != nil {
				return err
			}
		}

		IssuesList.Title = " Issues "
		StatusesList.SetTitle(" Projects > Statuses ")

		return nil
	})

	return nil
}

func Quit(g *ui.Gui, v *ui.View) error {
	return ui.ErrQuit
}
