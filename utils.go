package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
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

func LoadSites() error {
	ProjectsList.SetTitle("Sites")

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
func UpdateIssues() error {
	IssuesList.Reset()
	// Details.Clear()

	issues, _ := ListIssuesByProjectCode("FF")
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
