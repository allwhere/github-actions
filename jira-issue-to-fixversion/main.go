package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"regexp"
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
	issueKeys := extractIssueKeys(prTitle, projectKey)
	if len(issueKeys) == 0 {
		log.Println("No issue reference found in the PR title")
		return
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

func extractIssueKeys(prTitle string, projectKey string) []string {
	reBracket := regexp.MustCompile(`\[(.*?)\]`)
	bracketContents := reBracket.FindStringSubmatch(prTitle)

	if len(bracketContents) > 1 {
		reIssues := regexp.MustCompile(fmt.Sprintf(`%s-\d+`, projectKey))
		return reIssues.FindAllString(bracketContents[1], -1)
	}
	return []string{}
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
