package main

import (
	"fmt"
	"os"

	"github.com/ethanthoma/github-issue-data/pkg"
	issuequery "github.com/ethanthoma/github-issue-data/pkg/issue-query"
	"github.com/ethanthoma/github-issue-data/pkg/search"
)

func main() {
    token := os.Getenv("GITHUB_TOKEN")

    githubClient := github.NewClient(token)

    searchParams := search.NewSearch(
        search.Created("<=2019-09-30"),
        search.Is("public"),
        search.Fork(false),
        search.Mirror(false),
        search.Stars("0..1000"),
        search.Template(false),
    )

    perPage := 10
    page := 1

    repos, _, _, err := githubClient.FetchRepos(searchParams, perPage, page)
    if err != nil {

        fmt.Println("Error on fetching repos.\n[ERROR] -", err)
        panic(err)
    }

    issueQueryParams := issuequery.NewIssueQuery(
        issuequery.Labels("bug"),
        issuequery.State("closed"),
        issuequery.Page(1),
        issuequery.PerPage(10),
    )

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
