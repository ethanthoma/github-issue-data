package github

import (
	"encoding/json"
	"fmt"

	issuequery "github.com/ethanthoma/github-issue-data/pkg/issue-query"
	"github.com/ethanthoma/github-issue-data/pkg/repos"
)

func (client *Client) FetchRepos(fetchReposParams *repos.FetchReposParams) ([]Repo, int, bool, error) {
	url := "https://api.github.com/search/repositories"

	if fetchReposParams != nil {
		url += "?"
		if fetchReposParams.Search != nil {
			url += "q=" + fetchReposParams.Search.ToString()
		}

		url += fmt.Sprintf("&per_page=%d&page=%d", fetchReposParams.PerPage, fetchReposParams.Page)

		if fetchReposParams.Sort != nil {
			url += fmt.Sprintf("&sort=%s", fetchReposParams.Sort.Value)

			if fetchReposParams.Order != nil {
				url += fmt.Sprintf("&order=%s", fetchReposParams.Order.Value)
			}
		}
	}

	resp, err := client.fetch(url)
	if err != nil {
		return nil, 0, false, err
	}

	var result struct {
		TotalCount        int    `json:"total_count"`
		IncompleteResults bool   `json:"incomplete_results"`
		Items             []Repo `json:"items"`
	}

	err = json.Unmarshal(resp.Body, &result)

	return result.Items, result.TotalCount, result.IncompleteResults, err
}

func (client Client) FetchCommentsForIssue(repoFullname string, issueNumber int) ([]Comment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", repoFullname, issueNumber)
	resp, err := client.fetch(url)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	json.Unmarshal(resp.Body, &comments)
	return comments, nil
}

func (client *Client) UserIsCollaborator(repoFullname string, username string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/collaborators/%s", repoFullname, username)
	resp, err := client.fetch(url)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 204, nil
}

func (client *Client) FetchIssues(repoFullname string, issueQuery *issuequery.IssueQuery) ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues?%s", repoFullname, issueQuery.ToString())
	resp, err := client.fetch(url)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	json.Unmarshal(resp.Body, &issues)
	return issues, nil
}
