package main

import (
	"fmt"
    "os"
	"github.com/ethanthoma/github-issue-data/pkg"
)

func main() {
    token := os.Getenv("GITHUB_TOKEN")

    githubClient := github.NewClient(token)

    _ = githubClient

    searchParameters := map[string]string{
        "created":  "<=2019-09-30",
        "is":       "public",
        "fork":     "false",
        "mirror":   "false",
        "stars":    "0..1000",
        "template": "false",
    }
    perPage := 10
    page := 1

    repos, _, _, err := githubClient.FetchRepos(searchParameters, perPage, page)
    if err != nil {

        fmt.Println("Error on fetching repos.\n[ERROR] -", err)
        panic(err)
    }

    issueQueryParams := map[string]string{
        "labels":   "bug",
        "state":    "closed",
        "page":     "1",
        "per_page": "10",
    }

    for _, repo := range repos {
        issues, err := githubClient.FetchIssues(repo.FullName, issueQueryParams)
        if err != nil {
            fmt.Println("Error on fetching issues for ", repo.FullName, ".\n[ERROR] -", err)
            panic(err)
        }

        fmt.Println("Number of issues for", repo.FullName, "is", len(issues))


        for _, issue := range issues {
            fmt.Println("Number of comments in issue ID#", issue.ID, "is", issue.Comments)

            if issue.Comments == 0 {
                fmt.Println("Skipping...")
                continue
            }

            comments, err := githubClient.FetchCommentsForIssue(repo.FullName, issue.Number)
            if err != nil {
                fmt.Println("Error on fetching comments for issue.\n[ERROR] -", err)
                panic(err)
            }

            for _, comment := range comments {
                isCollab, err := githubClient.UserIsCollaborator(repo.FullName, comment.User.Login)
                if err != nil {
                    fmt.Println("Error on fetching user for comment.\n[ERROR] -", err)
                    panic(err)
                }
                fmt.Println("User:", comment.User.Login, "User is collab:", isCollab, "Text:", comment.Body)
            }
        }
    }
}
