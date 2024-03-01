package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
    "strings"
)

type Repo struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
}

type Issue = struct {
    ID          int     `json:"id"`
    URL         string  `json:"url"`
    Title       string  `json:"title"`
    Body        string  `json:"body"`
    User        User    `json:"user"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

type User struct {
    Login       string  `json:"login"`
    Type        string  `json:"type"`
    SiteAdmin   bool    `json:"site_admin"`
}

type Comment struct {
    ID          int     `json:"id"`
    Body        string  `json:"body"`
    User        User    `json:"user"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

var token = os.Getenv("GITHUB_TOKEN")

var headers = http.Header{
    "Accept": {"application/vnd.github+json"},
    "Authorization": {"Bearer "+token},
    "X-GitHub-Api-Version": {"2022-11-28"},
}

func fetch(url string) ([]byte, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error on request.\n[ERROR] -", err)
        return nil, err
    }
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error on response.\n[ERROR] -", err)
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

type SearchParameters = map[string]string

func fetchRepos(searchParameters SearchParameters, perPage int, page int) ([]Repo, int, bool, error) {
    var pairs []string
    for key, value := range searchParameters {
        pairs = append(pairs, key+":"+value)
    }
    params := strings.Join(pairs, "+")

    query := fmt.Sprintf("%s&sort=stars&order=desc&per_page=%d&page=%d", params, perPage, page)
    url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s", query)
    body, err := fetch(url)
    if err != nil {
        return nil, 0, false, err
    }

    var result struct {
        TotalCount        int    `json:"total_count"`
        IncompleteResults bool   `json:"incomplete_results"`
        Items             []Repo `json:"items"`
    }

    err = json.Unmarshal(body, &result)
    return result.Items, result.TotalCount, result.IncompleteResults, err
}

func fetchIssues(repoFullname string, queryParams map[string]string) ([]Issue, error) {
    var pairs []string
    for key, value := range queryParams {
        pairs = append(pairs, key+"="+value)
    }
    params := strings.Join(pairs, "&")

    url := fmt.Sprintf("https://api.github.com/repos/%s/issues?%s", repoFullname, params)
    body, err := fetch(url)
    if err != nil {
        return nil, err
    }

    var issues []Issue
    json.Unmarshal(body, &issues)
    return issues, nil
}

func fetchCommentsForIssue(repoFullname string, issueNumber int) ([]Comment, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", repoFullname, issueNumber)
    body, err := fetch(url)
    if err != nil {
        return nil, err
    }

    var comments []Comment
    json.Unmarshal(body, &comments)
    return comments, nil
}

func main() {
    searchParameters := SearchParameters{
        "created":  "<=2019-09-30",
        "is":       "public",
        "fork":     "false",
        "label":    "bug",
        "mirror":   "false",
        "stars":    "0..1000",
        "template": "false",
    }
    perPage := 10
    page := 1

    repos, _, _, err := fetchRepos(searchParameters, perPage, page)
    if err != nil {
        fmt.Println("Error on fetching repos.\n[ERROR] -", err)
        panic(err)
    }

    issueQueryParams := SearchParameters{
        "labels": "bug",
        "state": "closed",
        "page": "1",
        "per_page": "10",
    }

    for _, repo := range repos {
        issues, err := fetchIssues(repo.FullName, issueQueryParams)
        if err != nil {
            fmt.Println("Error on fetching issues for ", repo.FullName, ".\n[ERROR] -", err)
            panic(err)
        }

        fmt.Println("Number of issues for", repo.FullName, "is", len(issues))

        for _, issue := range issues {
            comments, err := fetchCommentsForIssue(repo.FullName, issue.ID)
            if err != nil {
                fmt.Println("Error on fetching comments for issue.\n[ERROR] -", err)
                panic(err)
            }

            fmt.Println("Number of comments in issue ID#", issue.ID, "is", len(comments))

            for _, comment := range comments {
                fmt.Println("User:", comment.User.Login, "Text:", comment.Body)
            }
        }
    }
}
