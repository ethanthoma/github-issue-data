package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github-issue-data/pkg"
	issuesquery "github-issue-data/pkg/issue"
)

type CommentData struct {
	RepoId    int    `json:"repo_id"`
	IssueID   int    `json:"issue_id"`
	CommentID int    `json:"comment_id"`
	AuthorID  int    `json:"author_id"`
	Author    string `json:"author"`
	Interval  int    `json:"interval"`
	Text      string `json:"text"`
	Type      string `json:"type"`
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("Please set the GITHUB_TOKEN environment variable.")
	}

	client := github.NewClient(token)

	sampleFilePath := "data/sample.csv"

	comments, err := getComments(client, sampleFilePath)
	if err != nil {
		fmt.Println("Error on getting comments.\n[ERROR] -", err)
		fmt.Print(client.RequestCount)
	}

	if comments != nil {
		fmt.Println("Number of comments:", len(*comments))
		github.SaveToCSV(comments, "data/comments.csv")
	}
}

func getComments(client *github.Client, sampleFilePath string) (*[]CommentData, error) {
	dataset := []CommentData{}

	file, err := os.Open(sampleFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	allRecords, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := allRecords[0]
	columnIndex := make(map[string]int)
	for i, columnName := range header {
		columnIndex[columnName] = i
	}

	allRecords = allRecords[1:]

	for i, record := range allRecords {
		parsedComments := len(dataset)

		repo := github.Repo{}
		for columnName, index := range columnIndex {
			switch columnName {
			case "id":
				repo.ID, _ = strconv.Atoi(record[index])
			case "name":
				repo.Name = record[index]
			case "full_name":
				repo.FullName = record[index]
			case "stargazers_count":
				repo.Stars, _ = strconv.Atoi(record[index])
			}
		}

		issues, err := filterIssues(client, &repo)
		if err != nil {
			fmt.Println("Failed to fetch issues for", repo.FullName, ":", err)
			return &dataset, err
		}

		dataset = append(dataset, (*issues)...)
		fmt.Println("Repos parsed:", i+1, "/", len(allRecords), "| Comments parsed:", len(dataset)-parsedComments)
	}

	return &dataset, nil
}

func filterIssues(client *github.Client, repo *github.Repo) (*[]CommentData, error) {
	data := []CommentData{}

	query := issuesquery.NewIssueQuery(
		issuesquery.State(issuesquery.Closed()),
		issuesquery.PerPage(100),
	)

	for page := 1; ; page++ {
		query.Set(issuesquery.Page(page))

		issues, err := client.FetchIssues(repo.FullName, query)
		if err != nil {
			fmt.Println("Failed to fetch issues for", repo.FullName, ":", err)
			return nil, err
		}

		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			if filterIssue(&issue) {
				comments, err := convertIssueToComments(client, repo, &issue)
				if err != nil {
					return nil, err
				}
				data = append(data, (*comments)...)
			}
		}
	}

	return &data, nil
}

func filterIssue(issue *github.Issue) bool {
	year := issue.CreatedAt.Year()
	return year > 2016 && year < 2020 && issue.PullRequest == nil && issue.State == "closed"
}

func convertIssueToComments(client *github.Client, repo *github.Repo, issue *github.Issue) (*[]CommentData, error) {
	data := []CommentData{}

	comments, err := client.FetchCommentsForIssue(repo.FullName, issue.Number)
	if err != nil {
		fmt.Println("Failed to fetch comments for", repo.FullName, ":", issue.ID)
		return nil, err
	}

	interval := dateToInterval(issue.CreatedAt)

	data = append(data, CommentData{
		RepoId:    repo.ID,
		IssueID:   issue.ID,
		CommentID: -1,
		AuthorID:  issue.User.ID,
		Author:    issue.User.Login,
		Interval:  interval,
		Text:      issue.Title + " " + issue.Body,
		Type:      issue.Type,
	})

	for _, comment := range comments {
		year := comment.CreatedAt.Year()
		if year > 2016 && year < 2020 {
			interval := dateToInterval(comment.CreatedAt)
			data = append(data, CommentData{
				RepoId:    repo.ID,
				IssueID:   issue.ID,
				CommentID: comment.ID,
				AuthorID:  comment.User.ID,
				Author:    comment.User.Login,
				Interval:  interval,
				Text:      comment.Body,
				Type:      comment.Type,
			})
		}
	}

	return &data, nil
}

func dateToInterval(date time.Time) int {
	startOfYear := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	weeks := int(date.Sub(startOfYear).Hours()/24/7) + 1
	return weeks
}
