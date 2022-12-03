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
		apiService := hackerNews.ApiHackerNewsClient{TopNewsEndpoint: config.HackerNewsGetTopNewsIds, GetByIdEndpoint: config.HackerNewsGetByItemid}
		news := apiService.GetNews(getExistNews(config.OldNewsFilePath).Items)
		for i := 0; i < len(news); i++ {
			if news[i].Url == "" {
				continue
			} else {
				bot.SendMessage(news[i].Title+"\n"+news[i].Url, config.TelegramChatId)

			}

		}
		writeNews(news, config.OldNewsFilePath)
		time.Sleep(200 * time.Second)
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

func getExistNews(filePath string) OldNews {
	var oldNews OldNews
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

func writeNews(news []*hackerNews.News, filePath string) {
	if news == nil {
		return
	}
	oldNews := getExistNews(filePath)
	var newsIds = oldNews.Items
	for i := range news {
		i := append(newsIds, news[i].Id)
		newsIds = i
	}
	writeNews := OldNews{newsIds}
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

type OldNews struct {
	Items []int64 `json:"items"`
}
