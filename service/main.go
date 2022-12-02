package main

import (
	"encoding/json"
	"fmt"
	hacker_news "hacker-news"
	"io"
	"log"
	"os"
	telegramBot "telegram-bot"
)

func main() {
	config := getConfiguration("../config.json")
	bot := telegramBot.TgBot{TelegramBotApiKey: config.TelegramBotApiKey}
	bot.StartBot()
	apiService := hacker_news.ApiServiceClient{ApiEndpoint: "https://hacker-news.firebaseio.com/v0/item/33805723.json?print=pretty"}
	fmt.Println(apiService.GetNews())
	bot.SendMessage(apiService.GetNews())

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

/*
Configuration App.
*/
type Configuration struct {
	TelegramBotApiKey       string `json:"telegramBotApiKey"`
	TelegramChatId          int64  `json:"telegramChatId"`
	HackerNewsGetByItemid   string `json:"hackerNewsGetByItemid"`
	HackerNewsGetTopNewsIds string `json:"hackerNewsGetTopNewsIds"`
}
