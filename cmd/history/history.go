package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethanthoma/github-issue-data/pkg"
)

type RepoHistory struct {
	RepoID   int `json:"repo_id"`
	Stars    int `json:"stars"`
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

	repoHistory, err := getRepoHistory(client, repos)

	if err != nil {
		fmt.Println("Error on getting history.")
		panic(err)
	}

	if repoHistory != nil {
		fmt.Println("Number of data fetched:", len(*repoHistory))
		github.SaveToCSV(repoHistory, "data/repoHistory-"+fmt.Sprint(len(*repoHistory))+".csv")
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
			fmt.Println("EOF: ", err)
			break
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

	for _, repo := range *repos {
		per_page := 100

		for page := 1; ; page++ {
			stars, err := client.FetchStarsForRepo(repo.FullName, per_page, page)
			if err != nil {
				fmt.Println("Error on fetching star history for repo:", repo.FullName)
				return nil, err
			}

			commits, err := client.FetchAllCommitsForRepo(repo.FullName)
			if err != nil {
				fmt.Printf("Error fetching commits for repo %s: %v\n", repo.FullName, err)
				return nil, err
			}

			if len(stars)+len(commits) == 0 {
				break
			}

			intervalData := make(map[int]*RepoHistory)

			for _, star := range stars {
				if star.StarredAt.Year() > 2016 && star.StarredAt.Year() < 2020 {
					continue
				}

				interval := dateToInterval(star.StarredAt)
				if _, ok := intervalData[interval]; !ok {
					intervalData[interval] = &RepoHistory{RepoID: repo.ID, Interval: interval}
				}
				intervalData[interval].Stars++
			}

			for _, commit := range commits {
				// choose the earlier of the two dates to determine the commit's interval
				commitDate := commit.Commit.Author.Date
				if commitDate.After(commit.Commit.Commiter.Date) {
					commitDate = commit.Commit.Commiter.Date
				}

				interval := dateToInterval(commitDate)
				if _, ok := intervalData[interval]; !ok {
					intervalData[interval] = &RepoHistory{RepoID: repo.ID, Interval: interval}
				}
				intervalData[interval].Commits++
			}

			for _, data := range intervalData {
				dataset = append(dataset, *data)
			}
		}
	}

	return &dataset, nil
}

func dateToInterval(date time.Time) int {
	startOfYear := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	weeks := int(date.Sub(startOfYear).Hours()/24/7) + 1
	return weeks
}
