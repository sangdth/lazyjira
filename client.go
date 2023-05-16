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
		log.Panicln(err)
	}

	return client, nil
}

func ListIssuesByProjectCode(projectCode string) ([]jira.Issue, error) {
	client, _ := GetJiraClient()

	// Define JQL query
	jql := fmt.Sprintf("project = %s", projectCode)

	// Get list of issues
	issues, _, err := client.Issue.Search(context.Background(), jql, nil)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
