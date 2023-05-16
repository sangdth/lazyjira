package main

import (
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

// func GetTicketsAssignedMe(jiraClient *jira.Client, ticketID string) ([]jira.Ticket, error) {
// }

// func GetTickets(jiraClient *jira.Client, projectCode string) (*jira.Ticket, error) {
// }
