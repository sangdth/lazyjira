package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	ui "github.com/awesome-gocui/gocui"
	config "github.com/gookit/config/v2"
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

// Creates a general purpose view to be used as input source
func createPromptView(g *ui.Gui, o CreateDialogOptions) {
	tw, th := g.Size()
	v, err := g.SetView(PromptView, tw/6, (th/2)-8, (tw*5)/6, (th/2)-6, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln(err)
	}

	g.Cursor = true

	PromptDialog = CreateDialog(v, PROMPT)
	PromptDialog.SetTitles(o.title, DialogDescription)
	PromptDialog.SetContent(o.content)
	PromptDialog.SetValue(o.value)
	PromptDialog.Focus(g)
}

func deletePromptView(g *ui.Gui) {
	g.Cursor = false
	if err := g.DeleteView(PromptView); err != nil {
		log.Panicln("Error while deleting prompt view", err)
	}
}

// Creates a view to be used as error alert or confirmation dialog
func createAlertView(g *ui.Gui, o CreateDialogOptions) {
	tw, th := g.Size()
	v, err := g.SetView(AlertView, tw/6, (th/2)-12, (tw*5)/6, (th/2)-6, 0)
	if err != nil && err != ui.ErrUnknownView {
		log.Panicln("Error while creating alert view", err)
	}

	g.Cursor = false

	AlertDialog = CreateDialog(v, ALERT)
	AlertDialog.SetTitles(o.title, DialogDescription)
	AlertDialog.SetContent(o.content)
	AlertDialog.SetValue(o.value)
	AlertDialog.Focus(g)
}

func deleteAlertView(g *ui.Gui) {
	g.Cursor = false
	if err := g.DeleteView(AlertView); err != nil {
		log.Panicln("Error while deleting alert view", err)
	}
}

func AddProject(g *ui.Gui, v *ui.View) error {
	ProjectsList.Unfocus()

	createPromptView(g, CreateDialogOptions{title: InsertNewCodeTitle})

	return nil
}

// Used when user press Esc to close or cancel a dialog
func CancelDialog(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case PromptView:
		if _, err := g.View(ProjectsView); err == nil && isNewCodeView(v) {
			ProjectsList.Focus(g)
		}
		if _, err := g.View(IssuesView); err == nil && isCreatingBranchView(v) {
			IssuesList.Focus(g)
		}

		deletePromptView(g)

		return nil

	case AlertView:
		deleteAlertView(g)
		if _, err := g.View(PromptView); err == nil {
			PromptDialog.Focus(g)
		} else {
			ProjectsList.Focus(g)
		}
		return nil
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
				existOpts := CreateDialogOptions{
					title:   " Alert! ",
					content: "Project already exist",
				}
				createAlertView(g, existOpts)

				return nil
			}

			statuses, _, err := SearchStatusesByProjectCode(value)
			if err != nil {
				errOpts := CreateDialogOptions{
					title:   " Alert! ",
					content: err.Error(),
				}
				createAlertView(g, errOpts)

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

		if isCreatingBranchView(v) {
			cmd := exec.Command("git", "branch", "-b", value)
			err := cmd.Run()
			if err != nil {
				log.Printf("------- %s", err)
			}

			log.Printf("Created branch: %s", value)
			deletePromptView(g)
			IssuesList.Focus(g)

			return nil
		}

		return nil
	})

	return nil
}

// Used when user press Enter to confirm a dialog (delete or ignore something)
func SubmitAlert(g *ui.Gui, v *ui.View) error {
	value := strings.TrimSpace(AlertDialog.value)
	if len(value) == 0 {
		return nil
	}

	g.Update(func(g *ui.Gui) error {
		if isDeleteView(v) {
			projectPath := fmt.Sprintf("%s.%s", ProjectsKey, strings.ToLower(value))
			if err := config.Set(projectPath, nil); err != nil {
				log.Panicln("Error while deleting project row", err)
			}
			writeConfigToFile()
			deleteAlertView(g)
			loadProjects()
			ProjectsList.Focus(g)

			return nil
		}

		return nil
	})

	return nil
}

// Move cursor up on list
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

// Move cursor down on list
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
		ProjectsList.Unfocus()
		if strings.Contains(IssuesList.Title, "bookmarks") {
			g.SelFgColor = ui.ColorMagenta | ui.AttrBold
		}
		IssuesList.Focus(g)
		return nil

	case StatusesView:
		StatusesList.Unfocus()
		IssuesList.Focus(g)
		return nil

	case IssuesView:
		if _, err := g.View(ProjectsView); err == nil {
			ProjectsList.Focus(g)
		}
		if _, err := g.View(StatusesView); err == nil {
			StatusesList.Focus(g)
		}

		IssuesList.Unfocus()
		return nil
	}

	return nil
}

// The Projects view has another second tab for Statuses
func SwitchProjectTab(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case StatusesView:
		ProjectsList.Focus(g)
		StatusesList.Unfocus()
		if err := g.DeleteView(StatusesView); err != nil {
			return err
		}
		return nil

	case ProjectsView:
		if err := createStatusView(g); err != nil {
			return err
		}
		ProjectsList.Unfocus()
		StatusesList.Focus(g)
		if err := OnEnterProject(g, v); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func ToggleStatus(g *ui.Gui, v *ui.View) error {
	currentItem := StatusesList.CurrentItem()
	if currentItem == "" {
		return nil
	}

	currentCursor := StatusesList.currentCursorY()
	projectCode := strings.ToLower(IssuesList.code)
	richStatusKey := strings.ToLower(currentItem)
	statusKey := richStatusKey[4:]

	path := fmt.Sprintf("%s.%s.statuses.%s", ProjectsKey, projectCode, statusKey)

	if config.Exists(path) {
		isChecked := config.Bool(path)
		if err := config.Set(path, !isChecked); err != nil {
			log.Panicln("Error while toggling status", err)
		}
	} else {
		if err := config.Set(path, true); err != nil {
			log.Panicln("Error while adding new status value", err)
		}
	}

	g.Update(func(g *ui.Gui) error {
		if err := FetchStatuses(g, projectCode); err != nil {
			StatusesList.SetTitle(" Projects > Statuses (Error!) ")
			return nil
		}

		if err := StatusesList.SetCursor(0, currentCursor); err != nil {
			return err
		}

		if err := FetchIssues(g, projectCode); err != nil {
			return err
		}

		return nil
	})

	return nil
}

// Pressing Spacebar will trigger this one
func OnSelectProject(g *ui.Gui, v *ui.View) error {
	projectCode := ProjectsList.CurrentItem()
	if projectCode == "" {
		return nil
	}

	IssuesList.Clear()

	IssuesList.SetTitle(" Issues | Fetching... ")
	IssuesList.SetCode(projectCode)

	// Can not nest the update
	g.Update(func(g *ui.Gui) error {
		if err := FetchIssues(g, projectCode); err != nil {
			return err
		}

		IssuesList.SetTitle(" Issues ")

		return nil
	})

	return nil
}

func OnEnterProject(g *ui.Gui, v *ui.View) error {
	currentItem := ProjectsList.CurrentItem()
	if currentItem == "" {
		return nil
	}

	projectCode := currentItem

	IssuesList.SetCode(projectCode)

	StatusesList.SetTitle(" Projects > Statuses | Fetching... ")
	IssuesList.SetTitle(" Issues | Fetching... ")

	isSameProject := IssuesList.code == projectCode

	g.Update(func(g *ui.Gui) error {
		if IssuesList.IsEmpty() || !isSameProject {
			if err := FetchIssues(g, projectCode); err != nil {
				IssuesList.SetTitle(" Issues (Error!) ")
				return nil
			}
		}

		if err := FetchStatuses(g, projectCode); err != nil {
			StatusesList.SetTitle(" Projects > Statuses (Error!) ")
			return nil
		}

		IssuesList.SetTitle(" Issues ")
		StatusesList.SetTitle(fmt.Sprintf(" Projects > Statuses (%s)", projectCode))
		ProjectsList.SetTitle(" Projects ")

		return nil
	})

	return nil
}

func RemoveProject(g *ui.Gui, v *ui.View) error {
	currentItem := ProjectsList.CurrentItem()
	if currentItem == "" {
		return nil
	}

	projectCode := currentItem

	ProjectsList.Unfocus()

	deleteOpts := CreateDialogOptions{
		title: DeleteConfirmTitle,
		content: fmt.Sprintf(`
			The project [%s] will be deleted.
			Action can not undo.
			Do you want to proceed?`, projectCode),
		value: projectCode,
	}
	createAlertView(g, deleteOpts)

	return nil
}

// Create git branch from selected issue
func GitBranchPrompt(g *ui.Gui, v *ui.View) error {
	currentItem := IssuesList.CurrentItem()
	if currentItem == "" {
		return nil
	}

	issueName := strings.ReplaceAll(currentItem, " ", "-")

	prefix := config.String(GitPrefixKey)

	branchName := issueName
	if prefix != "" {
		branchName = fmt.Sprintf("%s/%s", prefix, issueName)
	}

	createPromptView(g, CreateDialogOptions{
		title:   NewBranchTitle,
		content: branchName,
		value:   branchName,
	})

	if err := PromptDialog.SetCursor(len(branchName), 0); err != nil {
		return err
	}

	return nil
}

func Quit(g *ui.Gui, v *ui.View) error {
	writeConfigToFile()

	return ui.ErrQuit
}
