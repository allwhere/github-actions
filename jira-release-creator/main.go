package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"strconv"
)

var jiraClient *jira.Client
var jiraBaseUrl string

func main() {
	// Setup JIRA client
	if err := setupJiraClient(); err != nil {
		log.Fatalf("Error setting up Jira client: %v", err)
	}

	// Inputs from GitHub Actions
	versionName := os.Getenv("INPUT_VERSION_NAME")
	projectKey := os.Getenv("INPUT_PROJECT_KEY")

	// Get Project ID from Project Key
	projectID, err := getProjectID(projectKey)
	if err != nil {
		log.Fatalf("Error getting project ID: %v", err)
	}

	// Create Jira Version using the retrieved Project ID
	versionID, versionURL, err := createJiraVersion(versionName, projectID)
	if err != nil {
		log.Fatalf("Error creating Jira version: %v", err)
	}

	fmt.Printf("::set-output name=version-id::%s\n", versionID)
	fmt.Printf("::set-output name=version-url::%s\n", versionURL)
}

func getProjectID(projectKey string) (int, error) {
	req, err := jiraClient.NewRequest("GET", fmt.Sprintf("/rest/api/2/project/%s", projectKey), nil)
	if err != nil {
		return 0, fmt.Errorf("creating request failed: %v", err)
	}

	var project struct {
		ID string `json:"id"` // Assuming the ID is provided as a string in the JSON response
	}

	_, err = jiraClient.Do(req, &project)
	if err != nil {
		return 0, fmt.Errorf("request failed: %v", err)
	}

	// Convert project ID from string to integer
	projectID, convErr := strconv.Atoi(project.ID)
	if convErr != nil {
		return 0, fmt.Errorf("failed to convert project ID to integer: %v", convErr)
	}
	return projectID, nil
}

func setupJiraClient() error {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		return err
	}
	jiraClient = client
	return nil
}

func createJiraVersion(versionName string, projectID int) (string, string, error) {
	released := false
	version := jira.Version{
		Name:      versionName,
		ProjectID: projectID,
		Released:  &released,
	}

	createdVersion, resp, err := jiraClient.Version.Create(&version)
	if err != nil {
		return "", "", fmt.Errorf("failed to create Jira version: %s, error: %v", resp.Status, err)
	}

	versionURL := fmt.Sprintf("%v/projects/%s/versions/%s", os.Getenv("JIRA_URL"), projectID, createdVersion.ID)
	return createdVersion.ID, versionURL, nil
}
