package github

type Repo struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    FullName    string  `json:"full_name"`
}

type Issue = struct {
    ID          int     `json:"id"`
    URL         string  `json:"url"`
    Number      int     `json:"number"`
    Title       string  `json:"title"`
    Body        string  `json:"body"`
    User        User    `json:"user"`
    Comments    int     `json:"comments"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

type User struct {
    Login       string  `json:"login"`
    Type        string  `json:"type"`
}

type Comment struct {
    ID          int     `json:"id"`
    Body        string  `json:"body"`
    User        User    `json:"user"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}
