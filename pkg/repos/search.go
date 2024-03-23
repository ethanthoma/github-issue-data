package repos

import (
	"fmt"
	"strings"
)

type SearchParams struct {
	internal map[string]string
}

func NewSearchParams(options ...func(*SearchParams)) *SearchParams {
	searchParams := &SearchParams{
		internal: map[string]string{},
	}

	for _, option := range options {
		option(searchParams)
	}

	return searchParams
}

func (searchParams *SearchParams) Set(options ...func(*SearchParams)) {
	for _, option := range options {
		option(searchParams)
	}
}

func (from *SearchParams) Copy() *SearchParams {
	to := &SearchParams{
		internal: map[string]string{},
	}

	for k, v := range from.internal {
		to.internal[k] = v
	}

	return to
}

func setParam(key string, value string) func(*SearchParams) {
	return func(search *SearchParams) {
		search.internal[key] = value
	}
}

func Query(value string) func(*SearchParams) {
	return setParam("query", fmt.Sprint(value))
}

func Created(value Time) func(*SearchParams) {
	return setParam("created", fmt.Sprint(value.value))
}

func Pushed(value Time) func(*SearchParams) {
	return setParam("pushed", fmt.Sprint(value.value))
}

func Is(value IsValue) func(*SearchParams) {
	return setParam("is", fmt.Sprint(value.value))
}

func Fork(value bool) func(*SearchParams) {
	return setParam("fork", fmt.Sprint(value))
}

func Mirror(value bool) func(*SearchParams) {
	return setParam("mirror", fmt.Sprint(value))
}

func Stars(value Int) func(*SearchParams) {
	return setParam("stars", fmt.Sprint(value.value))
}

func Template(value bool) func(*SearchParams) {
	return setParam("template", fmt.Sprint(value))
}

func (searchParams SearchParams) ToString() string {
	var pairs []string
	query := ""
	for key, value := range searchParams.internal {
		if key == "query" {
			query = value + "+"
		} else {
			pairs = append(pairs, key+":"+value)
		}
	}
	return query + strings.Join(pairs, "+")
}

type FetchReposParams struct {
	Search  *SearchParams
	PerPage int
	Page    int
	Sort    *Sort
	Order   *Order
}

func NewFetchReposParams(options ...func(*FetchReposParams)) *FetchReposParams {
	fetchReposParams := &FetchReposParams{
		Search:  NewSearchParams(),
		PerPage: 30,
		Page:    1,
		Sort:    nil,
		Order:   nil,
	}

	for _, option := range options {
		option(fetchReposParams)
	}

	return fetchReposParams
}

func SetSearchParams(searchParams *SearchParams) func(*FetchReposParams) {
	return func(fetchReposParams *FetchReposParams) {
		fetchReposParams.Search = searchParams
	}
}

func SetPerPage(perPage int) func(*FetchReposParams) {
	return func(fetchReposParams *FetchReposParams) {
		fetchReposParams.PerPage = perPage
	}
}

func SetPage(page int) func(*FetchReposParams) {
	return func(fetchReposParams *FetchReposParams) {
		fetchReposParams.Page = page
	}
}

func SetSort(sort *Sort) func(*FetchReposParams) {
	return func(fetchReposParams *FetchReposParams) {
		fetchReposParams.Sort = sort
	}
}

func SetOrder(order *Order) func(*FetchReposParams) {
	return func(fetchReposParams *FetchReposParams) {
		fetchReposParams.Order = order
	}
}
