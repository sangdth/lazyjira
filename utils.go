package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	ui "github.com/awesome-gocui/gocui"
	color "github.com/gookit/color"
	config "github.com/gookit/config/v2"
	yaml "github.com/gookit/config/v2/yaml"
	keyring "github.com/zalando/go-keyring"
)

func getPaths() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	configDir := fmt.Sprintf("%s/%s", configHome, ProjectName)
	configPath := fmt.Sprintf("%s/%s", configDir, "config.yaml")

	return configPath
}

func initConfigSetup() {
	red := color.FgRed.Render

	configPath := getPaths()

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	if err := config.LoadFiles(configPath); err != nil {
		log.Fatalf("Missing config file, create one at %s", red(ConfigPathMsg))
	}
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
	projectCodeMap := config.StringMap(ProjectsKey)

	var projects []string
	for key := range projectCodeMap {
		projects = append(projects, strings.ToUpper(key))
	}

	return projects
}

func GetSavedStatusesByProjectCode(code string) []string {
	statusesPath := fmt.Sprintf("%s.%s.statuses", ProjectsKey, strings.ToLower(code))
	statusMap := config.StringMap(statusesPath)

	keys := make([]string, 0, len(statusMap))
	for k := range statusMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	statuses := make([]string, len(keys))
	for index, key := range keys {
		keyPath := fmt.Sprintf("%s.%s", statusesPath, key)

		richLabel := fmt.Sprintf("[ ] %s", strings.ToUpper(key))
		if config.Bool(keyPath) {
			richLabel = fmt.Sprintf("[v] %s", strings.ToUpper(key))
		}

		statuses[index] = richLabel
	}

	return statuses
}

func loadProjects() {
	ProjectsList.SetTitle(makeTabNames(ProjectsView))

	savedProjects := GetSavedProjects()

	sort.Strings(savedProjects)

	if len(savedProjects) == 0 {
		ProjectsList.SetTitle("No projects (Press 'a' to add)")
		ProjectsList.Reset()
		IssuesList.Reset()
		IssuesList.SetTitle("No issues")
	}

	ProjectsList.SetItems(savedProjects)
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

func FetchIssues(g *ui.Gui, code string) error {
	IssuesList.Reset()
	IssuesList.SetCode(code)

	issues, err := SearchIssuesByProjectCode(code)
	if err != nil {
		IssuesList.Title = fmt.Sprintf(" Failed to load issues from: %s ", code)
		IssuesList.Clear()
		return err
	}

	if len(issues) == 0 {
		IssuesList.SetTitle(fmt.Sprintf("No issues in %s", code))
		return nil
	}

	parsedIssues := make([]string, len(issues))
	for index, issue := range issues {
		key := issue.Key
		summary := issue.Fields.Summary
		row := fmt.Sprintf("%-2s %s", key, summary)
		parsedIssues[index] = row
	}

	IssuesList.SetItems(parsedIssues)

	return nil
}

func FetchStatuses(g *ui.Gui, code string) error {
	StatusesList.Reset()
	StatusesList.SetCode(code)

	oldParsedStatuses := GetSavedStatusesByProjectCode(code)

	if len(oldParsedStatuses) > 0 {
		StatusesList.SetItems(oldParsedStatuses)
		return nil
	}

	statuses, _, err := SearchStatusesByProjectCode(code)
	if err != nil {
		StatusesList.Title = " Projects > Statuses | Fetched failed "
		StatusesList.Clear()
		return err
	}

	if len(statuses) == 0 {
		return nil
	}

	newParsedStatuses := make([]string, len(statuses))
	for index, status := range statuses {
		newParsedStatuses[index] = fmt.Sprintf("[v] %s", strings.ToUpper(status.Name))
	}

	StatusesList.SetItems(newParsedStatuses)

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

	configPath := getPaths()

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

// func startSpinner(g *ui.Gui, v *ui.View) {
// 	spinnerInterval := 110 * time.Millisecond

// 	oldTitle := v.Title
// 	// Set the initial spinner state
// 	spinnerState := 0
// 	spinnerFrames := []string{"|", "/", "-", "\\"}

// 	// Create a function to update the view's title with the spinner
// 	updateTitle := func(v *ui.View) {
// 		v.Title = fmt.Sprintf(" %s %s ", oldTitle, spinnerFrames[spinnerState])
// 		spinnerState = (spinnerState + 1) % len(spinnerFrames)
// 	}

// 	go func() {
// 		for {
// 			g.Update(func(g *ui.Gui) error {
// 				if v, err := g.View(StatusesView); err == nil {
// 					updateTitle(v)
// 				}
// 				return nil
// 			})
// 			time.Sleep(spinnerInterval)
// 		}
// 	}()
// }

func isNewUsernameView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertUsernameTitle) // || strings.Contains(v.Title, "try again")
}

func isNewServerView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertServerTitle)
}

func isNewCodeView(v *ui.View) bool {
	return strings.Contains(v.Title, InsertNewCodeTitle)
}

func isDeleteView(v *ui.View) bool {
	return strings.Contains(v.Title, DeleteConfirmTitle)
}
