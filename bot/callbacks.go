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
		//! TODO: editing logic
		// TODO: получить запись из БД
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
	case data == "help":
	case data == "abort":
		// TODO: delete(userStates, chatID)
		// TODO: replace current keyboard with main keyboard
	}

	bot.Send(msg)
}
