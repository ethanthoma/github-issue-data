package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github-issue-data/pkg"

	"github.com/machinebox/graphql"
)

type StarHistory struct {
	RepoID   int `json:"repo_id"`
	Stars    int `json:"stars"`
	Interval int `json:"interval"`
}

type stargazer struct {
	StarredAt time.Time
}

func main() {
	reposFilepath := "data/sample.csv"
	repos, err := readRepos(reposFilepath)
	if err != nil {
		fmt.Println("Error fetching repos.")
		panic(err)
	}

	fmt.Println("Loaded sample repos.")

	client := graphql.NewClient("https://api.github.com/graphql")

	fmt.Println("Fetching stargazers.")

	var allStargazers []StarHistory
	historyParsed := 0
	for i, repo := range *repos {
		historyParsed = len(allStargazers)
		stargazers, err := fetchStargazers(client, repo)
		if err != nil {
			fmt.Println("Error fetching stargazers.", err)
		}
		allStargazers = append(allStargazers, stargazers...)
		fmt.Println("Repos parsed:", i+1, "/", len(*repos), "| Records added:", len(allStargazers)-historyParsed)
	}

	github.SaveToCSV(&allStargazers, "data/stargazers.csv")
}

func readRepos(filepath string) (*[]github.Repo, error) {
	repos := []github.Repo{}

	file, err := os.Open(filepath)
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

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading record from CSV:", err)
			return nil, err
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

		repos = append(repos, repo)
	}

	return &repos, nil
}

func fetchStargazers(client *graphql.Client, repo github.Repo) ([]StarHistory, error) {
	maxRetries := 50
	req := graphql.NewRequest(`
        query ($owner: String!, $name: String!, $cursor: String) {
            repository(owner: $owner, name: $name) {
                stargazers(first: 100, after: $cursor, orderBy: {field: STARRED_AT, direction: ASC}) {
                    edges {
                        starredAt
                        node {
                            login
                        }
                    }
                    pageInfo {
                        endCursor
                        hasNextPage
                    }
                }
            }
        }
	`)

	owner := strings.Split(repo.FullName, "/")[0]
	req.Var("owner", owner)
	req.Var("name", repo.Name)
	req.Var("cursor", nil)

	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Println("Please set GITHUB_TOKEN environment variable.")
		return nil, errors.New("missing GITHUB_TOKEN")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN")))

	var respData struct {
		Repository struct {
			Stargazers struct {
				Edges []struct {
					StarredAt time.Time `json:"starredAt"`
					Node      struct {
						Login string `json:"login"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"stargazers"`
		} `json:"repository"`
	}

	stargazers := []stargazer{}

	cursor := ""
	for {
		req.Var("cursor", cursor)
		if err := fetchWithRetry(client, req, &respData, maxRetries); err != nil {
			fmt.Println("Failed to fetch stargazers after retries:", err)
			return nil, err
		}
		for _, edge := range respData.Repository.Stargazers.Edges {
			stargazers = append(stargazers, stargazer{StarredAt: edge.StarredAt})
		}
		if !respData.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}
		cursor = respData.Repository.Stargazers.PageInfo.EndCursor
	}

	if len(stargazers) == 0 {
		return nil, errors.New("no stargazers found")
	}

	sort.Slice(stargazers, func(i, j int) bool {
		return stargazers[i].StarredAt.Before(stargazers[j].StarredAt)
	})

	var stargazerHistories []StarHistory

	lastInterval := 0
	totalStars := 0
	for _, sg := range stargazers {
		if sg.StarredAt.Year() < 2016 {
			continue
		} else if sg.StarredAt.Year() > 2019 {
			break
		}

		interval := dateToInterval(sg.StarredAt)
		if interval != lastInterval {
			// skip if no stars accumulated yet
			if totalStars > 0 {
				stargazerHistories = append(stargazerHistories, StarHistory{
					RepoID:   repo.ID,
					Stars:    totalStars,
					Interval: lastInterval,
				})
			}
			lastInterval = interval
		}
		totalStars++
	}

	return stargazerHistories, nil
}

func fetchWithRetry(client *graphql.Client, req *graphql.Request, respData interface{}, maxRetries int) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = client.Run(context.Background(), req, respData)
		if err != nil {
			fmt.Printf("Attempt %d: request failed: %s\n", attempt, err)
			backoffDuration := time.Duration(attempt^2) * time.Second
			fmt.Printf("Waiting for %s before retrying...\n", backoffDuration)
			time.Sleep(backoffDuration)
			continue
		}
		return nil
	}
	return err
}

func dateToInterval(date time.Time) int {
	startOfYear := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	weeks := int(date.Sub(startOfYear).Hours()/24/7) + 1
	return weeks
}
