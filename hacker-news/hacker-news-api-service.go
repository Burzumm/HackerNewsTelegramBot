package hacker_news

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ApiClient interface {
	GetTopAndNewNews() []int64
	GetNews(id string) *[]News
}

func (t ApiHackerNewsClient) GetNews(existNews []int64) []*News {
	var news []*News
	topAndNewNews := t.GetTopAndNewNews()
	var notExistNews []int64
	if len(existNews) != 0 {
		for newsIndex := 0; newsIndex < len(topAndNewNews); newsIndex++ {
			isAdded := false
			for oldNewsIndex := 0; oldNewsIndex < len(existNews); oldNewsIndex++ {
				if topAndNewNews[newsIndex] == existNews[oldNewsIndex] {
					isAdded = true
				}
			}
			if !isAdded {
				i := append(notExistNews, topAndNewNews[newsIndex])
				notExistNews = i
			}
		}
	}

	for i := 0; i < len(notExistNews); i++ {
		url := t.GetByIdEndpoint + strconv.FormatInt(notExistNews[i], 10) + ".json?print=pretty"
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)
		newsList := append(news, newsCreate(res))
		err = res.Body.Close()
		if err != nil {
			return nil
		}
		news = newsList
	}

	return news

}
func (t ApiHackerNewsClient) GetTopAndNewNews() []int64 {
	payload := strings.NewReader("{}")
	req, _ := http.NewRequest("GET", t.TopNewsEndpoint, payload)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	return newsIdsCreate(res)
}

func newsIdsCreate(r *http.Response) []int64 {
	var news []int64
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		panic(err)
	}
	return news
}

func newsCreate(r *http.Response) *News {
	var news *News
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		panic(err)
	}
	return news
}

type ApiHackerNewsClient struct {
	TopNewsEndpoint string
	GetByIdEndpoint string
}

type News struct {
	By          string `json:"by"`
	Descendants int64  `json:"descendants"`
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Type        string `json:"type"`
	Time        int64  `json:"time"`
	Score       int64  `json:"score"`
	Text        string `json:"text"`
}
