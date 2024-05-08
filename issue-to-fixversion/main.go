package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"regexp"
	"strings"
)

var jiraClient *jira.Client

func main() {
	// Setup JIRA client
	if err := setupJiraClient(); err != nil {
		log.Fatalf("Error setting up Jira client: %v", err)
	}

	// Inputs from GitHub Actions
	prTitle := os.Getenv("PR_TITLE")
	projectKey := os.Getenv("PROJECT_KEY")
	fixVersion := os.Getenv("FIX_VERSION")

	// Extract Jira issue keys from PR title
	issueKeys, err := extractIssueKeys(prTitle, projectKey)
	if err != nil {
		log.Fatalf("Error extracting issue keys: %v", err)
	}
	fmt.Println("Extracted issue keys:", issueKeys)

	// Add issues to Jira fix version
	if err := addIssuesToFixVersion(issueKeys, projectKey, fixVersion); err != nil {
		log.Fatalf("Error adding issues to fix version: %v", err)
	}

	fmt.Println("Successfully added issues to fix version.")
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

func extractIssueKeys(prTitle, projectKey string) ([]string, error) {
	rePattern := fmt.Sprintf(`\[(%s-\d+(?:,%s-\d+)*)\]`, projectKey, projectKey)
	re := regexp.MustCompile(rePattern)
	matches := re.FindStringSubmatch(prTitle)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no issue keys found in PR title")
	}

	issueKeys := strings.Split(matches[1], ",")
	for i, key := range issueKeys {
		issueKeys[i] = strings.TrimSpace(key)
	}
	return issueKeys, nil
}

func addIssuesToFixVersion(issueKeys []string, projectKey string, fixVersion string) error {
	for _, issueKey := range issueKeys {
		issue, _, err := jiraClient.Issue.Get(issueKey, nil)
		if err != nil {
			return fmt.Errorf("failed to retrieve issue %s: %v", issueKey, err)
		}

		fixVersion := jira.FixVersion{Name: fixVersion}
		issue.Fields.FixVersions = append(issue.Fields.FixVersions, &fixVersion)

		_, _, err = jiraClient.Issue.Update(issue)
		if err != nil {
			return fmt.Errorf("failed to update issue %s with fix version %s: %v", issueKey, fixVersion.Name, err)
		}
	}
	return nil
}
