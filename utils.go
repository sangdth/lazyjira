package main

import (
	"bytes"
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
	PROJECT_NAME = "lazyjira"
	PROJECTS     = "projects"
	SERVER       = "server"
	USERNAME     = "username"
	HELP_LINK    = "https://github.com/sangdth/lazyjira#getting-started"
	JIRA_LINK    = "https://id.atlassian.com/manage-profile/security/api-tokens"
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
	configPath := fmt.Sprintf("%s/%s", home, PROJECT_NAME)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // Config file not found
			if err := os.Mkdir(configPath, 0755); err != nil {
				log.Panicln("Error while creating creating folder", err)
			}

			viper.SetConfigFile(fmt.Sprintf("%s/config.yaml", configPath))

			err := viper.WriteConfig()
			if err != nil {
				log.Panicln("Error while initiating config file", err)
			}

			InitConfigValue() // can NOT work, why???
			// panic: runtime error: invalid memory address or nil pointer dereference
			// [signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x1030fa664
		}
	}

	if !viper.IsSet(PROJECTS) {
		savedProjects := map[string]map[string]interface{}{
			ASSIGNED_TO_ME: {
				"statuses": nil,
			},
		}

		// TODO Still dont know how to reload after the first initial
		viper.Set(PROJECTS, savedProjects)

		if err := viper.WriteConfig(); err != nil {
			log.Panicln("Error while setting default saved projects", err)
		}
	}

	viper.WatchConfig()

	return nil
}

func GetJiraCredentials() (string, string, string, error) {
	server := viper.GetString(SERVER)
	username := viper.GetString(USERNAME)

	secret, err := keyring.Get(PROJECT_NAME, username)
	if err != nil {
		log.Panicln(err)
	}

	return server, username, secret, err
}

func GetSavedProjects() []string {
	var projects []string

	stringMap := viper.GetStringMap(PROJECTS)

	for key := range stringMap {
		projects = append(projects, strings.ToUpper(key))
	}

	return projects
}

func GetSavedStatusesByProjectCode(code string) []interface{} {
	stringMap := viper.GetStringMap(fmt.Sprintf("%s.%s.statuses", PROJECTS, code))

	statuses := make([]interface{}, len(stringMap))

	index := 0
	for key := range stringMap {
		statuses[index] = key
		index++
	}

	return statuses
}

func SetNewStatusesByProjectCode(code string, value []interface{}) error {
	path := fmt.Sprintf("%s.%s.statuses", PROJECTS, code)

	newValue := make(map[string]interface{}, len(value))
	for _, status := range value {
		newValue[status.(string)] = true
	}

	viper.Set(path, newValue)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func LoadProjects() {
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

	if len(issues) == 0 {
		IssuesList.SetTitle(fmt.Sprintf("No issues in %s", "FF"))
		return nil
	}

	data := make([]interface{}, len(issues))
	for index, issue := range issues {
		key := issue.Key
		summary := issue.Fields.Summary
		row := fmt.Sprintf("%-2s %s", key, summary)
		data[index] = row
	}

	IssuesList.SetItems(data)

	return nil
}

func RenderStatusesList(statuses []jira.Status) error {
	StatusesList.Reset()

	if len(statuses) == 0 {
		return nil
	}

	data := make([]interface{}, len(statuses))
	for index, status := range statuses {
		data[index] = strings.ToLower(status.Name)
	}

	StatusesList.SetItems(data)

	if err := SetNewStatusesByProjectCode(StatusesList.code, data); err != nil {
		return err
	}

	return nil
}

func FetchIssues(g *ui.Gui, code string) error {
	IssuesList.SetCode(code)

	issues, err := SearchIssuesByProjectCode(code)

	if err != nil {
		IssuesList.Title = fmt.Sprintf(" Failed to load issues from: %s ", code)
		IssuesList.Clear()
		return err
	}

	if err := RenderIssuesList(issues); err != nil {
		return err
	}

	return nil
}

func FetchStatuses(g *ui.Gui, code string) error {
	StatusesList.SetCode(code)

	statuses, _, err := SearchStatusesByProjectCode(code)
	if err != nil {
		StatusesList.Title = " Projects > Statuses | Fetched failed "
		StatusesList.Clear()

		// if err := createAlertView(g, "New Alert:"); err != nil {
		// 	log.Panicln("Error on createAlertView", err)
		// }

		return nil
	}

	if err := RenderStatusesList(statuses); err != nil {
		log.Println("Error on RenderStatuses")
		return err
	}

	return nil
}

func spaces(n int) string {
	var s bytes.Buffer
	for i := 0; i < n; i++ {
		s.WriteString(" ")
	}
	return s.String()
}
