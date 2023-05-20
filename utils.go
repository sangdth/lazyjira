package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	// fsnotify "github.com/fsnotify/fsnotify"
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

func InitConfig() {
	home, _ := GetConfigHome()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("%s/%s", home, projectName))

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}

	if !viper.IsSet("savedProjects") {
		savedProjects := map[string]map[string]interface{}{
			ASSIGNED_TO_ME: {
				"statuses": nil,
			},
		}

		// TODO Still dont know how to reload after the first initial
		viper.Set("savedProjects", savedProjects)

		err := viper.WriteConfig()
		if err != nil {
			log.Panicln(err)
		}
	}

	viper.WatchConfig()
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
	var projects []string

	stringMap := viper.GetStringMap("savedProjects")

	for key := range stringMap {
		projects = append(projects, strings.ToUpper(key))
	}

	return projects
}

func LoadProjects(v *ui.View) {
	ProjectsList.SetTitle(makeTabNames(ProjectsView))

	savedProjects := GetSavedProjects()

	if len(savedProjects) == 0 {
		ProjectsList.SetTitle("No projects (Ctrl-f to add)")
		ProjectsList.Reset()
		IssuesList.Reset()
		IssuesList.SetTitle("No issues")
	}
	data := make([]interface{}, len(savedProjects))
	for index, project := range savedProjects {
		data[index] = project
	}

	ProjectsList.SetItems(data)
}

func makeTabNames(name string) string {
	switch name {

	case ProjectsView:
		return " Projects "

	case StatusesView:
		return " Projects > Statuses "
	}

	return "Something went wrong in making name"
}

// Make the from for project key, currently hardcoded "FF"
func RenderIssuesList(issues []jira.Issue) error {
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

	IssuesList.SetItems(data)

	return nil
}

func RenderStatusesList(issues []jira.Issue) error {
	StatusesList.Reset()

	if len(issues) == 0 {
		//  TODO: We need a better way to handle tab title
		// StatusesList.SetTitle(fmt.Sprintf("No issues in %v", "FF"))
		return nil
	}

	statusesMap := make(map[string]bool)
	for _, issue := range issues {
		statusesMap[issue.Fields.Status.Name] = true
	}

	index := 0
	data := make([]interface{}, len(statusesMap))
	for status := range statusesMap {
		row := fmt.Sprintf(" %s", status) //  <-- for unchecked
		data[index] = row
		index++
	}

	StatusesList.SetItems(data)

	return nil
}

func FetchIssues(g *ui.Gui, code string) {
	IssuesList.SetCode(code)
	IssuesList.Title = " Issues | Fetching... "

	g.Update(func(g *ui.Gui) error {
		issues, err := SearchIssuesByProjectCode(code)
		if err != nil {
			IssuesList.Title = fmt.Sprintf(" Failed to load issues from: %v ", code)
			IssuesList.Clear()
		}

		if err := RenderIssuesList(issues); err != nil {
			log.Println("Error on RenderIssues", err)
			return err
		}

		return nil
	})
}

func FetchStatuses(g *ui.Gui, code string) {
	StatusesList.Title = " Projects > Statuses | Fetching... "
	g.Update(func(g *ui.Gui) error {
		issues, err := SearchIssuesByProjectCode(code)
		if err != nil {
			StatusesList.Title = " Projects > Statuses | Fetched failed "
			StatusesList.Clear()
		}

		if err := RenderStatusesList(issues); err != nil {
			log.Println("Error on RenderStatuses", err)
			return err
		}

		return nil
	})
}
