package main

import (
	"log"
	// "context"
	// "errors"
	// "fmt"
	// "regexp"
	// "strings"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
)

func GetJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: "",
		APIToken: "",
	}

	client, err := jira.NewClient("https://lokalise.atlassian.net", tp.Client())

	if err != nil {
		log.Panicln(err)
	}

	return client
}

// func GetTicketsAssignedMe(jiraClient *jira.Client, ticketID string) ([]jira.Ticket, error) {
// }

// func GetTickets(jiraClient *jira.Client, projectCode string) (*jira.Ticket, error) {
// }
