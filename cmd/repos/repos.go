package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ethanthoma/github-issue-data/pkg"
	"github.com/ethanthoma/github-issue-data/pkg/repos"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("Please set the GITHUB_TOKEN environment variable.")
	}

	client := github.NewClient(token)

	repos, err := getRepos(client)
	if err != nil {
		fmt.Println("Error on getting batch.\n[ERROR] -", err)
		panic(err)
	}
	github.SaveToCSV(repos, "data/repos.csv")
}

func getRepos(client *github.Client) (*[]github.Repo, error) {
	defaultSearchParams := getSearchFilter()

	populationSize, err := getPopulationSize(client, defaultSearchParams.Copy())
	if err != nil {
		fmt.Println("Error on getting population size.\n[ERROR] -", err)
		panic(err)
	}

	populationIds := make(map[int]bool)
	population := make([]github.Repo, populationSize)

	perPage := 100
	minStars := 100
	defaultSearchParams.Set(repos.Stars(repos.Int{}.Min(minStars)))

	then := time.Now()
	for page, index := 0, 0; index < populationSize; page++ {
		if page == 10 {
			page = 0
			minStars = population[index-1].Stars
			defaultSearchParams.Set(repos.Stars(repos.Int{}.Min(minStars)))
			fmt.Println("Progress: ", index+1, "/", populationSize)
			fmt.Println("minStars: ", minStars)
		}

		now := time.Now()
		diff := now.Sub(then)
		if diff.Seconds() < 2 {
			time.Sleep(diff)
		}

		repos, _, incomplete, err := client.FetchRepos(repos.NewFetchReposParams(
			repos.SetSearchParams(defaultSearchParams),
			repos.SetPage(page),
			repos.SetPerPage(perPage),
			repos.SetSort(repos.SortByStars()),
			repos.SetOrder(repos.Asc()),
		))
		if err != nil {
			fmt.Println("Error on fetching batch.\n[ERROR] -", err)
			return nil, err
		}
		then = now

		if incomplete {
			fmt.Println("[WARNING] incomplete page.")
		}

		for _, repo := range repos {
			if populationIds[repo.ID] {
				continue
			}

			population[index] = repo
			populationIds[repo.ID] = true
			index++
		}
		break
	}

	fmt.Println("Progress: ", populationSize, "/", populationSize)

	return &population, nil
}

func getSearchFilter() *repos.SearchParams {
	created, _ := time.Parse("2006-01-02", "2019-09-30")
	pushed := created.AddDate(0, 6, 0)
	searchFilter := repos.NewSearchParams(
		repos.Query("library"),
		repos.Created(repos.Time{}.Max(created)),
		repos.Is(repos.Public()),
		repos.Mirror(false),
		repos.Template(false),
		// API returns incorrect max if min is less than 20
		repos.Stars(repos.Int{}.Min(100)),
		repos.Pushed(repos.Time{}.Min(pushed)),
	)

	return searchFilter
}

func getPopulationSize(client *github.Client, searchParams *repos.SearchParams) (int, error) {
	_, populationSize, _, err := client.FetchRepos(repos.NewFetchReposParams(
		repos.SetSearchParams(searchParams),
		repos.SetPage(1),
		repos.SetPerPage(1),
		repos.SetSort(repos.SortByStars()),
	))

	if err != nil {
		fmt.Println("Error on getting population size.\n[ERROR] -", err)
		return 0, err
	}

	fmt.Println("population size:", populationSize)

	return populationSize, nil
}
