package github

import (
	"encoding/json"
	"fmt"
	"time"

	issuequery "github-issue-data/pkg/issue"
	"github-issue-data/pkg/repos"
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

func (client *Client) FetchIssues(repoFullname string, issueQuery *issuequery.IssueQuery) ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues?%s", repoFullname, issueQuery.ToString())
	resp, err := client.fetch(url)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	err = json.Unmarshal(resp.Body, &issues)

	if err != nil {
		print(string(resp.Body))
		return nil, err
	}

	return issues, nil
}

func (client Client) FetchCommentsForIssue(repoFullname string, issueNumber int) ([]Comment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", repoFullname, issueNumber)
	resp, err := client.fetch(url)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	err = json.Unmarshal(resp.Body, &comments)

	if err != nil {
		print(string(resp.Body))
		return nil, err
	}

	return comments, nil
}

func (client *Client) FetchAllCommitsForRepo(repoFullname string, since time.Time, until time.Time, perPage, page int) ([]Commit, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/commits?since=%s&until=%sper_page=%d&page=%d",
		repoFullname,
		since.Format("2006-01-02T15:04:05Z"),
		until.Format("2006-01-02T15:04:05Z"),
		perPage,
		page,
	)
	resp, err := client.fetch(url)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	err = json.Unmarshal(resp.Body, &commits)

	if err != nil {
		print(string(resp.Body))
		return nil, err
	}

	return commits, nil
}
