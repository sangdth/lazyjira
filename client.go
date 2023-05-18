package main

import (
	"context"
	"fmt"
	"log"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
)

func GetJiraClient() (*jira.Client, error) {
	// Get Jira API token from Keychain
	server, username, secret, _ := GetJiraCredentials()

	tp := jira.BasicAuthTransport{
		Username: username,
		APIToken: secret,
	}

	client, err := jira.NewClient(server, tp.Client())

	if err != nil {
		log.Panicln("Failed to initiate new Jira client", err)
	}

	return client, nil
}

func MakeJQL(code string) string {
	switch code {

	case "Assigned to me":
		return "assignee=currentUser()"

	default:
		return fmt.Sprintf("project=%s", code)
	}
}

func ListIssuesByProjectCode(projectCode string) ([]jira.Issue, error) {
	client, _ := GetJiraClient()

	// Define JQL query
	jql := MakeJQL(projectCode)

	// Get list of issues
	issues, _, err := client.Issue.Search(context.Background(), jql, nil)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
