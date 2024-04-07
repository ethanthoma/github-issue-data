package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github-issue-data/pkg"
)

type RepoHistory struct {
	RepoID   int `json:"repo_id"`
	Commits  int `json:"commits"`
	Interval int `json:"interval"`
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("Please set the GITHUB_TOKEN environment variable.")
	}

	client := github.NewClient(token)

	reposFilepath := "data/sample.csv"
	repos, err := readRepos(reposFilepath)
	if err != nil {
		fmt.Println("Error fetching repos.")
		panic(err)
	}

	fmt.Println("Loaded sample repos.")

	repoHistory, err := getRepoHistory(client, repos)

	if err != nil {
		fmt.Println("Error on getting history.")
		panic(err)
	}

	if repoHistory != nil {
		fmt.Println("Number of records fetched:", len(*repoHistory))
		github.SaveToCSV(repoHistory, "data/commits.csv")
	}
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

func getRepoHistory(client *github.Client, repos *[]github.Repo) (*[]RepoHistory, error) {
	dataset := []RepoHistory{}

	since := time.Date(2016, 01, 01, 0, 0, 0, 0, time.UTC)
	until := time.Date(2019, 12, 31, 23, 59, 59, 9999, time.UTC)

	addedRecords := 0
	for i, repo := range *repos {
		addedRecords = len(dataset)
		per_page := 100

		intervalData := make(map[int]*RepoHistory)

		for page := 1; ; page++ {
			commits, err := client.FetchAllCommitsForRepo(repo.FullName, since, until, per_page, page)
			if err != nil {
				fmt.Printf("Error fetching commits for repo %s: %v\n", repo.FullName, err)
				return nil, err
			}

			if len(commits) == 0 {
				break
			}

			for _, commit := range commits {
				var commitDate time.Time

				if commit.Commit.Author.Date.Year() == 0 {
					commitDate = commit.Commit.Commiter.Date
				} else {
					commitDate = commit.Commit.Author.Date
				}

				interval := dateToInterval(commitDate)
				if _, ok := intervalData[interval]; !ok {
					intervalData[interval] = &RepoHistory{RepoID: repo.ID, Interval: interval}
				}
				intervalData[interval].Commits++
			}
		}

		for _, data := range intervalData {
			dataset = append(dataset, *data)
		}

		fmt.Println("Repos parsed:", i+1, "/", len(*repos), "| Records added:", len(dataset)-addedRecords)
	}

	return &dataset, nil
}

func dateToInterval(date time.Time) int {
	startOfYear := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	weeks := int(date.Sub(startOfYear).Hours()/24/7) + 1
	return weeks
}
