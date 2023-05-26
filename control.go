package main

import (
	"fmt"
	"log"
	"strings"

	ui "github.com/awesome-gocui/gocui"
	config "github.com/gookit/config/v2"
)

var (
	PromptDialog *Dialog
	AlertDialog  *Dialog
)

func createStatusView(g *ui.Gui) error {
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

// createPromptView creates a general purpose view to be used as input source
// from the user
func createPromptView(g *ui.Gui, title string) error {
	tw, th := g.Size()
	v, err := g.SetView(PromptView, tw/6, (th/2)-8, (tw*5)/6, (th/2)-6, 0)
	if err != nil && err != ui.ErrUnknownView {
		return err
	}

	g.Cursor = true

	PromptDialog = CreateDialog(v, PROMPT)
	PromptDialog.SetTitles(title, " (Press Esc to close) ")
	PromptDialog.Focus(g)

	return nil
}

// deletePromptView deletes the current prompt view
func deletePromptView(g *ui.Gui) {
	g.Cursor = false
	if err := g.DeleteView(PromptView); err != nil {
		log.Panicln("Error while deleting prompt view", err)
	}
}

func createAlertView(g *ui.Gui, msg string) {
	tw, th := g.Size()
	v, err := g.SetView(AlertView, tw/6, (th/2)-12, (tw*5)/6, (th/2)-6, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Error while creating alert view", err)
	}

	g.Cursor = false

	AlertDialog = CreateDialog(v, ALERT)
	AlertDialog.SetTitles(" Error! ", " (Press Esc to close) ")
	AlertDialog.SetContent(msg)
	AlertDialog.Focus(g)
}

func deleteAlertView(g *ui.Gui) error {
	g.Cursor = false
	return g.DeleteView(AlertView)
}

func AddProject(g *ui.Gui, v *ui.View) error {
	ProjectsList.Unfocus()

	if err := createPromptView(g, InsertNewCodeTitle); err != nil {
		log.Panicln("Error on AddProject", err)
	}

	return nil
}

func CloseFloatView(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case PromptView:
		deletePromptView(g)
		ProjectsList.Focus(g)

	case AlertView:
		if err := deleteAlertView(g); err != nil {
			log.Println("Error on deletePromptView", err)
			return err
		}
		PromptDialog.Focus(g)
		g.Cursor = true

	}
	return nil
}

func SubmitPrompt(g *ui.Gui, v *ui.View) error {
	value := strings.TrimSpace(v.ViewBuffer())
	if len(value) == 0 {
		return nil
	}

	g.Update(func(g *ui.Gui) error {
		if isNewUsernameView(v) {
			if err := config.Set(UsernameKey, value); err != nil {
				log.Panicln("Error while init username", err)
			}
			writeConfigToFile()
			deletePromptView(g)
			loadProjects()
			ProjectsList.Focus(g)

			return nil
		}

		if isNewServerView(v) {
			if err := config.Set(ServerKey, value); err != nil {
				log.Panicln("Error while init server", err)
			}
			writeConfigToFile()
			deletePromptView(g)
			loadProjects()
			ProjectsList.Focus(g)

			return nil
		}

		if isNewCodeView(v) {
			path := fmt.Sprintf("%s.%s", ProjectsKey, value)

			if config.Exists(path) {
				createAlertView(g, "Project already exist")

				return nil
			}

			statuses, _, err := SearchStatusesByProjectCode(value)
			if err != nil {
				createAlertView(g, err.Error())

				return nil
			}

			convertedStatuses := make(map[string]bool, len(statuses))
			for _, status := range statuses {
				convertedStatuses[strings.ToLower(status.Name)] = true
			}

			newValue := map[string]map[string]bool{
				"statuses": convertedStatuses,
			}
			if err := config.Set(path, newValue); err != nil {
				log.Panicln("Error while setting new statuses", err)
			}

			writeConfigToFile()
			deletePromptView(g)
			loadProjects()
			ProjectsList.Focus(g)

			return nil
		}

		return nil
	})

	return nil
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
			log.Println("Error on SitesList", err)
			return err
		}
	case StatusesView:
		if err := StatusesList.MoveDown(); err != nil {
			log.Println("Error on StatusesList", err)
			return err
		}
	case IssuesView:
		if err := IssuesList.MoveDown(); err != nil {
			log.Println("Error on IssuesList", err)
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
		if err := createStatusView(g); err == nil {
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

	projectCode := strings.ToLower(IssuesList.code)

	statusKey := strings.ToLower(currentItem.(string))

	path := fmt.Sprintf("%s.%s.statuses.%s", ProjectsKey, projectCode, statusKey)

	if config.Exists(path) {
		currentValue := config.Bool(path)
		if err := config.Set(path, !currentValue); err != nil {
			log.Panicln("Error while toggling status", err)
		}
	} else {
		if err := config.Set(path, true); err != nil {
			log.Panicln("Error while adding new status value", err)
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

	if err := createStatusView(g); err != nil {
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
	writeConfigToFile()

	return ui.ErrQuit
}
