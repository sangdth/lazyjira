package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	ui "github.com/awesome-gocui/gocui"
	viper "github.com/spf13/viper"
	keyring "github.com/zalando/go-keyring"
)

const (
	projectName      = "lazyjira"
	helpLink         = "https://github.com/sangdth/lazyjira#getting-started"
	jiraAPITokenLink = "https://id.atlassian.com/manage-profile/security/api-tokens"
)

func GetConfigHome() (string, error) {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home != "" {
		return home, nil
	}
	return home + "/.config", nil
}

func InitConfig() error {
	home, _ := GetConfigHome()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("%s/%s", home, projectName))

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}

	return nil
}

func GetJiraCredentials() (string, string, string, error) {
	server := viper.GetString("server")
	username := viper.GetString("login")

	secret, err := keyring.Get(projectName, username)
	if err != nil {
		log.Panicln(err)
	}

	return server, username, secret, err
}

func GetSavedProjects() []string {
	projects := viper.GetStringSlice("savedProjects")

	return projects
}

func LoadProjects(v *ui.View) error {
	ProjectsList.SetTitle(MakeProjectTabNames(ProjectsView))

	savedProjects := GetSavedProjects()

	if len(savedProjects) == 0 {
		ProjectsList.SetTitle("No projects (Ctrl-f to add)")
		ProjectsList.Reset()
		IssuesList.Reset()
		IssuesList.SetTitle("No issues")
		return nil
	}
	data := make([]interface{}, len(savedProjects))
	for index, project := range savedProjects {
		data[index] = project
	}

	return ProjectsList.SetItems(data)
}

// Make the from for project key, currently hardcoded "FF"
func UpdateIssuesList(issues []jira.Issue) error {
	IssuesList.Reset()
	// Details.Clear()

	if len(issues) == 0 {
		IssuesList.SetTitle(fmt.Sprintf("No issues in %v", "FF"))
		return nil
	}
	IssuesList.SetTitle(fmt.Sprintf("Issues from: %v", "FF"))

	data := make([]interface{}, len(issues))
	for index, issue := range issues {
		// if _, ok := eventInBookmarks(e); ok {
		// 	e.Title = fmt.Sprintf("  %v", e.Title)
		// }
		key := issue.Key
		summary := issue.Fields.Summary
		row := fmt.Sprintf("%-2s %s", key, summary)
		data[index] = row
	}

	return IssuesList.SetItems(data)
}

func UpdateStatusesList(issues []jira.Issue) error {
	StatusesList.Reset()

	if len(issues) == 0 {
		// StatusesList.SetTitle(fmt.Sprintf("No issues in %v", "FF"))
		return nil
	}
	// IssuesList.SetTitle(fmt.Sprintf("Issues from: %v", "FF"))

	columnStatuses := make(map[string]bool)

	data := make([]interface{}, len(issues))
	for index, issue := range issues {
		value := issue.Fields.Status.Name
		// key := issue.Key
		// row := fmt.Sprintf("%-2s %s", key, value)
		data[index] = value
		columnStatuses[issue.Fields.Status.Name] = true
		// log.Println(issue.Fields.Status.Name)
	}

	return StatusesList.SetItems(data)
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

	g.SelFgColor = ui.ColorGreen | ui.AttrBold

	return nil
}

func SwitchProjectTab(g *ui.Gui, v *ui.View) error {
	switch v.Name() {

	case StatusesView:
		ProjectsList.Focus(g)
		StatusesList.Unfocus()
		g.DeleteView(StatusesView)

	case ProjectsView:
		if err := createStatusView(g); err == nil {
			OnEnter(g, v)
			ProjectsList.Unfocus()
			StatusesList.Focus(g)
		} else {
			log.Panicln("Error on createStatusView()", err)
		}
	}

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

func FetchIssues(g *ui.Gui, code string) error {
	IssuesList.Title = " Issues | Fetching... "
	g.Update(func(g *ui.Gui) error {
		issues, err := ListIssuesByProjectCode(code)
		if err != nil {
			IssuesList.Title = fmt.Sprintf(" Failed to load issues from: %v ", code)
			IssuesList.Clear()
		}

		if err := UpdateIssuesList(issues); err != nil {
			log.Println("Error on UpdateIssues", err)
			return err
		}

		if err := UpdateStatusesList(issues); err != nil {
			log.Println("Error on UpdateStatuses", err)
			return err
		}

		return nil
	})

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
	err := OnSelectProject(g, v)

	return err
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
