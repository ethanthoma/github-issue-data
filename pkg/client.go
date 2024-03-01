package github

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
    httpClient      *http.Client 
    headers         http.Header
}

func NewClient(token string) *Client {
    return &Client {
        httpClient: &http.Client{},
        headers: http.Header{
            "Accept": {"application/vnd.github+json"},
            "Authorization": {"Bearer "+token},
            "X-GitHub-Api-Version": {"2022-11-28"},
        },
    }
}

type Response struct {
    StatusCode  int
    Body        []byte
}

func (client *Client) fetch(url string) (*Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error on request.\n[ERROR] -", err)
        return nil, err
    }
    req.Header = client.headers

    resp, err := client.httpClient.Do(req)
    if err != nil {
        fmt.Println("Error on response.\n[ERROR] -", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    return &Response{
        resp.StatusCode,
        body,
    }, err
}
