package bot

import (
	"fmt"

	"github.com/afoninartem/delivery-price/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainMenuKB() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonURL("Открыть мини-апп", "https://deliver-price-mini-app.vercel.app/"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить локацию", "new_loc"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать локации", "edit_loc"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("\U0001F4B0 Узнать цены \U0001F4B0", "get_prices"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
		),
	)
	return keyboard
}

func editLocationsKB(locs []models.Location) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, loc := range locs {
		nameBtn := tgbotapi.NewInlineKeyboardButtonData(loc.Name, "_")
		renameBtn := tgbotapi.NewInlineKeyboardButtonData("\U0001F4DD", fmt.Sprintf("rnm:%d", loc.ID))
		deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑️", fmt.Sprintf("del:%d", loc.ID))
		row := tgbotapi.NewInlineKeyboardRow(nameBtn, renameBtn, deleteBtn)
		keyboard = append(keyboard, row)
	}
	mmRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩️ В главное меню", "abort"))
	keyboard = append(keyboard, mmRow)
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
