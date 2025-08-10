package bot

import (
	"log/slog"
	"os"

	"github.com/afoninartem/delivery-price/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botApiKey      = os.Getenv("botApiKey")
	bot            *tgbotapi.BotAPI
	userStates     = make(map[int64]*models.UserState)
	userLastPrices = make(map[int64]map[uint]string)
)

func Bot() {
	b, err := tgbotapi.NewBotAPI(botApiKey)
	if err != nil {
		slog.Error("can't start bot with key botApiKey", "error", err)
	}
	bot = b

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		// case update.Message != nil && update.Message.IsCommand():
		// 	handleCommand(update.Message.Command(), update.Message.Chat.ID)
		case update.CallbackQuery != nil:
			handleCallback(update.CallbackQuery)
		default:
			if update.Message != nil {
				if _, e := userStates[update.Message.Chat.ID]; !e {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите один из пунктов меню.")
					msg.ReplyMarkup = mainMenuKB()
					bot.Send(msg)
					continue
				}
				handleConversation(update)
			}

		}
	}
}
