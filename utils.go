package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	ui "github.com/awesome-gocui/gocui"
	config "github.com/gookit/config/v2"
	yaml "github.com/gookit/config/v2/yaml"
	keyring "github.com/zalando/go-keyring"
)

func getPaths() (string, string, string) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	configDir := fmt.Sprintf("%s/%s", configHome, ProjectName)
	configPath := fmt.Sprintf("%s/%s", configDir, "config.yaml")

	return configPath, configDir, configHome
}

func InitConfig() error {
	configPath, configDir, _ := getPaths()

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	// Can not find the folder, start creating it
	if _, err := os.Stat(configDir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			log.Printf("Error while creating config folder %s", err)
		}

		// Assume we don't have the file as well, so create it
		file, err := os.Create(configPath)
		if err != nil {
			log.Printf("Error while creating config file %s", err)
		}

		defer file.Close()

		if err := config.LoadFiles(configPath); err != nil {
			log.Panicln("Error while loading config file", err)
		}
	}

	// Can not find the folder, start creating it
	// if _, err := os.Stat(configDir); err != nil && os.IsNotExist(err) {
	// if err := os.Mkdir(configDir, 0755); err != nil {
	// 	log.Printf("Error while creating config folder %s", err)
	// }

	// // Assume we don't have the file as well, so create it
	// file, err := os.Create(configPath)
	// if err != nil {
	// 	log.Printf("Error while creating config file %s", err)
	// }

	// defer file.Close()

	// if !config.Exists(ProjectsKey) {
	// 	initData := map[string]map[string]bool{
	// 		"statuses": {
	// 			"open": true,
	// 		},
	// 	}

	// 	if err := config.Set(fmt.Sprintf("%s.%s", ProjectsKey, AssignedToMeKey), initData); err != nil {
	// 		log.Panicln("Error while init first configs", err)
	// 	}

	// 	// it is nil wtf????
	// 	log.Println("before: ", config.Get(fmt.Sprintf("%s.%s", ProjectsKey, AssignedToMeKey)))
	// 	writeConfigToFile()

	// 	return nil // need to remove this return
	// }

	// TODO Without this LoadFiles everything will be nil
	// how to create the file before run the init?
	if err := config.LoadFiles(configPath); err != nil {
		log.Panicln("Error while loading config file", err)
	}

	return nil
}

func GetJiraCredentials() (string, string, string, error) {
	server := config.String(ServerKey)
	username := config.String(UsernameKey)

	secret, err := keyring.Get(ProjectName, username)
	if err != nil {
		log.Panicln(err)
	}

	return server, username, secret, err
}

func GetSavedProjects() []string {
	stringMap := config.StringMap(ProjectsKey)

	var projects []string
	for key := range stringMap {
		projects = append(projects, strings.ToUpper(key))
	}

	return projects
}

func GetSavedStatusesByProjectCode(code string) []interface{} {
	stringMap := config.StringMap(fmt.Sprintf("%s.%s.statuses", ProjectsKey, code))

	statuses := make([]interface{}, len(stringMap))

	index := 0
	for key := range stringMap {
		statuses[index] = key
		index++
	}

	return statuses
}

func LoadProjects() {
	ProjectsList.SetTitle(makeTabNames(ProjectsView))

	savedProjects := GetSavedProjects()

	if len(savedProjects) == 0 {
		ProjectsList.SetTitle("No projects (Press 'a' to add)")
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

		return nil
	}

	if err := RenderStatusesList(statuses); err != nil {
		log.Println("Error on RenderStatuses")
		return err
	}

	return nil
}

/*
 * This helper will set the config into memory and write it to the file
 */
func writeConfigToFile() {
	buff := new(bytes.Buffer)

	if _, err := config.DumpTo(buff, config.Yaml); err != nil {
		log.Printf("Error while dumping config file: %s\n", err)
	}

	configPath, _, _ := getPaths()

	if err := os.WriteFile(configPath, buff.Bytes(), 0755); err != nil {
		log.Printf("Error while writing config file: %s\n", err)
	}
}

func spaces(n int) string {
	var s bytes.Buffer
	for i := 0; i < n; i++ {
		s.WriteString(" ")
	}
	return s.String()
}

func isNewUsernameView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertUsernameTitle) // || strings.Contains(v.Title, "try again")
}

func isNewServerView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertServerTitle)
}

func isNewCodeView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertNewCodeTitle)
}
