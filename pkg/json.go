package github

import "time"

type Repo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Stars    int    `json:"stargazers_count"`
}

type Issue = struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	User        User      `json:"user"`
	State       string    `json:"state"`
	Comments    int       `json:"comments"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PullRequest *struct{} `json:"pull_request,omitempty"`
	Type        string    `json:"author_association"`
}

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type Comment struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"author_association"`
}

type Star struct {
	StarredAt time.Time `json:"starred_at"`
}

type Commit struct {
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
		Commiter struct {
			Date time.Time `json:"date"`
		} `json:"commiter"`
	} `json:"commit"`
}
