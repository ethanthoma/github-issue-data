package main

import (
	"encoding/csv"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/ethanthoma/github-issue-data/pkg"
	issuesquery "github.com/ethanthoma/github-issue-data/pkg/issue"
)

type CommentData struct {
	IssueID   int
	CommentID int
	AuthorID  int
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string
	Type      string
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	client := github.NewClient(token)

	reposFilepath := "data/repos.csv"

	comments, err := getComments(client, reposFilepath, 500, 13_611, rand.NewPCG(420, 69))
	if err != nil {
		fmt.Println("Error on getting comments.\n[ERROR] -", err)
	}

	if comments != nil {
		fmt.Println("Number of comments:", len(*comments))
		github.SaveToCSV(comments, "data/comments-"+fmt.Sprint(len(*comments))+".csv")
	}
}

func getComments(client *github.Client, reposFilepath string, sampleSize int, populationSize int, seed *rand.PCG) (*[]CommentData, error) {
	dataset := []CommentData{}

	indices := *getIndices(sampleSize, populationSize, seed)

	file, err := os.Open(reposFilepath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV header:", err)
		return nil, err
	}

	columnIndex := make(map[string]int)
	for i, columnName := range header {
		columnIndex[columnName] = i
	}

	lines := 0
	for i, index := range indices {
		parsedComments := len(dataset)
		for ; ; lines++ {
			record, err := reader.Read()
			if err != nil {
				fmt.Println("EOF: ", err)
				break
			}

			if index < lines {
				continue
			}

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
			break
		}

		fmt.Println("Repos parsed:", i+1, "/", sampleSize, "| Comments parsed:", len(dataset)-parsedComments)
	}

	return &dataset, nil
}

func getIndices(sampleSize int, populationSize int, seed *rand.PCG) *[]int {
	random := rand.New(seed)

	indices := make([]int, sampleSize)
	generated := map[int]bool{}

	count := 0
	for count < sampleSize {
		index := random.IntN(populationSize)
		if !generated[index] {
			indices[count] = index
			generated[index] = true
			count++
		}
	}

	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})

	return &indices
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
			year := issue.CreatedAt.Year()
			if year >= 2016 && year <= 2019 && issue.PullRequest == nil {
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

func convertIssueToComments(client *github.Client, repo *github.Repo, issue *github.Issue) (*[]CommentData, error) {
	data := []CommentData{}

	comments, err := client.FetchCommentsForIssue(repo.FullName, issue.Number)
	if err != nil {
		fmt.Println("Failed to fetch comments for", repo.FullName, ":", issue.ID)
		return nil, err
	}

	data = append(data, CommentData{
		IssueID:   issue.ID,
		CommentID: -1,
		AuthorID:  issue.User.ID,
		Author:    issue.User.Login,
		CreatedAt: issue.CreatedAt,
		UpdatedAt: issue.UpdatedAt,
		Text:      issue.Title + " " + issue.Body,
		Type:      issue.Type,
	})

	for _, comment := range comments {
		year := comment.CreatedAt.Year()
		if year >= 2016 && year <= 2019 {
			data = append(data, CommentData{
				IssueID:   issue.ID,
				CommentID: comment.ID,
				AuthorID:  comment.User.ID,
				Author:    comment.User.Login,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
				Text:      comment.Body,
				Type:      comment.Type,
			})
		}
	}

	return &data, nil
}
