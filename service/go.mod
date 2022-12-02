module service

go 1.19

replace (
	hacker-news v0.0.0 => ../hacker-news
	telegram-bot v0.0.0 => ../telegram-bot
)

require telegram-bot v0.0.0

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
