package github

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	headers    http.Header
	limiter    *rate.Limiter
}

func NewClient(token string) *Client {
	limiter := rate.NewLimiter(rate.Limit(5000./(60.*60.)+0.01), 1)

	return &Client{
		httpClient: &http.Client{},
		headers: http.Header{
			"Accept":               {"application/vnd.github+json"},
			"Authorization":        {"Bearer " + token},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
		limiter: limiter,
	}
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (client *Client) fetch(url string) (*Response, error) {
	if err := client.limiter.Wait(context.Background()); err != nil {
		fmt.Println("Rate limiter error:", err)
		return nil, err
	}

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
		StatusCode: resp.StatusCode,
		Body:       body,
	}, err
}
