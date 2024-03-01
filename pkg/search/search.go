package search

import (
    "fmt"
    "strings"
)

type Search struct { 
    internal map[string]string
}

func NewSearch(options ...func(*Search)) *Search {
    searchParams := &Search{
        internal: map[string]string{},
    }

    for _, option := range options {
        option(searchParams)
    }

    return searchParams
}

func setParam(key string, value string) func(*Search) {
    return func(search *Search) {
        search.internal[key] = value
    }
}

func Created(value string) func(*Search) {
    return setParam("created", fmt.Sprint(value))
}

func Is(value string) func(*Search) {
    return setParam("is", fmt.Sprint(value))
}

func Fork(value bool) func(*Search) {
    return setParam("fork", fmt.Sprint(value))
}

func Mirror(value bool) func(*Search) {
    return setParam("mirror", fmt.Sprint(value))
}

func Stars(value string) func(*Search) {
    return setParam("stars", fmt.Sprint(value))
}

func Template(value bool) func(*Search) {
    return setParam("template", fmt.Sprint(value))
}

func (searchParams Search) ToString() string {
    var pairs []string
    for key, value := range searchParams.internal {
        pairs = append(pairs, key+":"+value)
    }
    return strings.Join(pairs, "+")
}

