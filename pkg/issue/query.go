package issue

import (
	"fmt"
	"strings"
)

type IssueQuery struct {
	internal map[string]string
}

func NewIssueQuery(options ...func(*IssueQuery)) *IssueQuery {
	issueQuery := &IssueQuery{
		internal: map[string]string{},
	}

	for _, option := range options {
		option(issueQuery)
	}

	return issueQuery
}

func (query *IssueQuery) Set(options ...func(*IssueQuery)) {
	for _, option := range options {
		option(query)
	}
}

func setParam(key string, value string) func(*IssueQuery) {
	return func(issueQuery *IssueQuery) {
		issueQuery.internal[key] = value
	}
}

func Labels(value string) func(*IssueQuery) {
	return setParam("labels", fmt.Sprint(value))
}

func Page(value int) func(*IssueQuery) {
	return setParam("page", fmt.Sprint(value))
}

func PerPage(value int) func(*IssueQuery) {
	return setParam("per_page", fmt.Sprint(value))
}

type QueryState struct {
	value string
}

func State(value QueryState) func(*IssueQuery) {
	return setParam("state", fmt.Sprint(value.value))
}

func Open() QueryState {
	return QueryState{"open"}
}

func Closed() QueryState {
	return QueryState{"closed"}
}

func All() QueryState {
	return QueryState{"all"}
}

func (issueQuery IssueQuery) ToString() string {
	var pairs []string
	for key, value := range issueQuery.internal {
		pairs = append(pairs, key+"="+value)
	}
	return strings.Join(pairs, "&")
}
