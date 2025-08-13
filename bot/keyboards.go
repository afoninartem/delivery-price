package bot

import (
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"

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
		nameBtn := tgbotapi.NewInlineKeyboardButtonData(loc.Name, "nocb")
		renameBtn := tgbotapi.NewInlineKeyboardButtonData("\U0001F4DD", fmt.Sprintf("rnm:%d", loc.ID))
		deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑️", fmt.Sprintf("del:%d", loc.ID))
		row := tgbotapi.NewInlineKeyboardRow(nameBtn, renameBtn, deleteBtn)
		keyboard = append(keyboard, row)
	}
	mmRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩️ В главное меню", "abort"))
	keyboard = append(keyboard, mmRow)
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func pricesKB(chatID int64, locs []models.Location) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	locsWithPrice := getPrices(chatID, locs)
	for _, loc := range locsWithPrice {
		priceRow := priceBtnRow(loc)
		row := tgbotapi.NewInlineKeyboardRow(priceRow...)
		keyboard = append(keyboard, row)
	}
	refresh := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "get_prices"))
	keyboard = append(keyboard, refresh)
	mmRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩️ В главное меню", "abort"))
	keyboard = append(keyboard, mmRow)
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func priceBtnRow(loc models.Location) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	nameBtn := tgbotapi.NewInlineKeyboardButtonData(loc.Name, "nocb")
	priceDif := priceDif(loc)
	row = append(row, nameBtn, priceDif)
	return row
}

func priceDif(loc models.Location) tgbotapi.InlineKeyboardButton {
	var text, emoji string
	if len(loc.LastPrice) == 0 {
		text = loc.Price
	} else {
		lp, err := decimal.NewFromString(loc.LastPrice)
		if err != nil {
			slog.Error("decimal from string", "error", err)
			lp = decimal.NewFromInt(90)
		}
		ap, err := decimal.NewFromString(loc.Price)
		if err != nil {
			slog.Error("decimal from string", "error", err)
		}
		diff := lp.Sub(ap)
		switch {
		case diff.GreaterThan(decimal.Zero):
			emoji = "🔴"
		case diff.LessThan(decimal.Zero):
			emoji = "🟢"
		default:
			emoji = "🟡"
		}
		text = fmt.Sprintf("%s%s%s", lp, emoji, ap)
	}

	return tgbotapi.NewInlineKeyboardButtonData(text, "nocb")
}
