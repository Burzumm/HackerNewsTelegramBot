package main

import (
	"encoding/json"
	hackerNews "hacker-news"
	"io"
	"io/ioutil"
	"log"
	"os"
	telegramBot "telegram-bot"
	"time"
)

func main() {
	config := getConfiguration("./config.json")
	bot := telegramBot.TgBot{TelegramBotApiKey: config.TelegramBotApiKey}
	bot.StartBot()
	for {
		var newsSend []hackerNews.News
		apiService := hackerNews.ApiHackerNewsClient{TopNewsEndpoint: config.HackerNewsGetTopNewsIds, GetByIdEndpoint: config.HackerNewsGetByItemid}
		existItems := getExistNews(config.OldNewsFilePath).Items
		news := apiService.GetNews(&existItems)
		for _, item := range *news {
			if item.Url == "" || item.Score < 100 {
				continue
			} else {
				bot.SendMessage(item.Title+"\n\n"+item.Url, config.TelegramChatId)
				newsSend = append(newsSend, item)

			}
		}
		writeNews(newsSend, config.OldNewsFilePath)
		time.Sleep(600 * time.Second)
	}

}

func getConfiguration(path string) *Configuration {
	var m Configuration
	jsonFile, errorFile := os.Open(path)
	if errorFile != nil {
		log.Panic(errorFile)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	errJson := json.Unmarshal(byteValue, &m)
	if errJson != nil {
		log.Panic(errJson)
	}
	return &m
}

func getExistNews(filePath string) OldNewsList {
	var oldNews OldNewsList
	jsonFile, errorFile := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if errorFile != nil {
		log.Printf("error open file : %s to : %s\n", filePath, errorFile)
		return oldNews
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Printf("error close file : %s to : %s\n", filePath, errorFile)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	errJson := json.Unmarshal(byteValue, &oldNews)
	if errJson != nil {
		log.Printf("error Unmarshal to : %s\n", errorFile)
		return oldNews
	}
	return oldNews
}

func writeNews(news []hackerNews.News, filePath string) {
	if news == nil {
		return
	}
	oldNews := getExistNews(filePath)
	var newsIds = oldNews.Items
	if len(newsIds) >= 10000 {
		newsIds = newsIds[len(newsIds)-500 : len(newsIds)-1]
	}
	for _, item := range news {
		newList := append(newsIds, item)
		newsIds = newList
	}
	writeNews := OldNewsList{newsIds}
	file, err := json.MarshalIndent(writeNews, "", " ")
	err = ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		log.Printf("error write to file : %s to : %s\n", filePath, err)
		return
	}
}

/*
Configuration App.
*/
type Configuration struct {
	TelegramBotApiKey       string `json:"telegramBotApiKey"`
	TelegramChatId          int64  `json:"telegramChatId"`
	HackerNewsGetByItemid   string `json:"hackerNewsGetByItemid"`
	HackerNewsGetTopNewsIds string `json:"hackerNewsGetTopNewsIds"`
	OldNewsFilePath         string `json:"oldNewsFilePath"`
}

type OldNewsList struct {
	Items []hackerNews.News `json:"items"`
}
