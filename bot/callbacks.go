package bot

import (
	"fmt"
	"strings"

	"github.com/afoninartem/delivery-price/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(cb *tgbotapi.CallbackQuery) {
	var (
		chatID = cb.Message.Chat.ID
		msgID  = cb.Message.MessageID
		data   = cb.Data
		msg    = tgbotapi.NewMessage(chatID, "")
	)

	switch {
	case data == "new_loc":
		msg.Text = "Введите скопированные с карт координаты."
		if _, e := userStates[chatID]; e {
			userStates[chatID].Step = "new_coords"
		} else {
			userStates[chatID] = &models.UserState{
				Step: "new_coords",
			}
		}

	case data == "edit_loc":
		locs, err := getUserLocs(chatID)
		if err != nil {
			msg.Text = "Не удалось получить список локаций, попробуйте позже."
			break
		}
		msg.Text = "Ваши сохраненные локации:"
		msg.ReplyMarkup = editLocationsKB(locs)
	case strings.Contains(data, "rnm:"):
		id := getID(data)
		loc, err := getLocByID(id)
		if err != nil {
			msg.Text = "Не удалось переименовать локацию, попробуйте позже."
			break
		}
		state := &models.UserState{
			Step:     "rename",
			Location: loc,
		}
		userStates[chatID] = state
		msg.Text = fmt.Sprintf("Введите новое название для локации %s.", loc.Name)
	case strings.Contains(data, "del:"):
		id := getID(data)
		loc, err := getLocByID(id)
		if err != nil {
			msg.Text = "Не удалось удалить локацию, попробуйте позже."
			break
		}
		err = loc.Delete()
		if err != nil {
			msg.Text = "Не удалось удалить локацию, попробуйте позже."
			break
		}
		msg.Text = fmt.Sprintf("Локация %s успешно удалена.", loc.Name)
		msg.ReplyMarkup = mainMenuKB()
	case data == "get_prices":
		callback := tgbotapi.NewCallback(cb.ID, "Поиск цен")
		bot.Send(callback)
		locs, err := getUserLocs(chatID)
		if err != nil {
			msg.Text = "Не удалось получить список локаций, попробуйте позже."
			break
		}
		msg.Text = "Актуальные расценки на тариф Экспресс:"
		msg.ReplyMarkup = pricesKB(chatID, locs)
		delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		bot.Send(delMsg)
	case data == "help":
		msg.Text = help()
	case data == "abort":
		delete(userStates, chatID)
		text := "Выберите действие."
		edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, mainMenuKB())
		bot.Send(edit)
		return
	case data == "nocb":
		ans := tgbotapi.NewCallback(cb.ID, "Нажатие на эту кнопку ни к чему не приведет.")
		ans.ShowAlert = false
		bot.Request(ans)
	}

	bot.Send(msg)
}
