package hacker_news

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type ApiService interface {
	GetNews() string
}

func (t ApiServiceClient) GetNews() string {
	res, err := http.Get(t.ApiEndpoint)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	return string(resBody)
}

type ApiServiceClient struct {
	ApiEndpoint string
}
