package main

import (
	"encoding/csv"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"

	"github.com/ethanthoma/github-issue-data/pkg"
	issuesquery "github.com/ethanthoma/github-issue-data/pkg/issue"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	client := github.NewClient(token)

	reposFilepath := "data/repos.csv"

	issues, err := getIssues(client, reposFilepath, 500, 13_611, rand.NewPCG(420, 69))
	if err != nil {
		fmt.Println("Error on getting issues.\n[ERROR] -", err)
		panic(err)
	}

	fmt.Println("Number of issues:", len(*issues))

	github.SaveToCSV(issues, "data/issues.csv")
}

func getIssues(client *github.Client, reposFilepath string, sampleSize int, populationSize int, seed *rand.PCG) (*[]github.Issue, error) {
	var dataset []github.Issue

	indices := *getIndices(sampleSize, populationSize, seed)

	fmt.Println(indices)

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
		fmt.Println("Repos parsed:", i, "/", sampleSize)

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
				return nil, err
			}

			dataset = append(dataset, (*issues)...)
			break
		}
	}

	fmt.Println("Repos parsed:", sampleSize, "/", sampleSize)

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

func filterIssues(client *github.Client, repo *github.Repo) (*[]github.Issue, error) {
	var filteredIssues []github.Issue

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
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	return &filteredIssues, nil
}
