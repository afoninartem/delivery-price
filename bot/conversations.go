package bot

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/afoninartem/delivery-price/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleConversation(u tgbotapi.Update) {
	var (
		chatID = u.Message.Chat.ID
		input  = u.Message.Text
		state  = userStates[chatID]
		msg    = tgbotapi.NewMessage(chatID, "")
	)
	switch state.Step {
	case "new_coords":
		if validateCoords(input) {
			text := "Принято! Теперь придумайте название для этого места."
			state.Location = handleCoords(input)
			state.Step = "new_name"
			userStates[chatID] = state
			msg.Text = text
		} else {
			text := "Введены невалидные координаты. Скопиуйте координаты нужной локации из Яндекс Карт и вставьте их сюда."
			msg.Text = text
		}
	case "new_name":
		state.Location.Name = input
		state.Location.UserID = chatID
		///// TODO: добавить локу в БД.
		err := state.Location.Create()
		if err != nil {
			msg.Text = err.Error()
			break
		}
		///// TODO: сообщить об этом юзеру.
		msg.Text = fmt.Sprintf("Локация %s с координатами %s, %s добавлена в Избранное.", state.Location.Name, state.Location.Lat, state.Location.Lng)
		msg.ReplyMarkup = mainMenuKB()
		delete(userStates, chatID)
	case "rename":
		loc := state.Location
		oldName := loc.Name
		loc.Name = input
		err := loc.Update()
		if err != nil {
			msg.Text = "Не удалось переименовать локацию, попробуйте позже."
			break
		}
		msg.Text = fmt.Sprintf("Локация %s переименована в %s.", oldName, loc.Name)
		msg.ReplyMarkup = mainMenuKB()
	}
	bot.Send(msg)
}

func validateCoords(input string) bool {
	coordRegex := regexp.MustCompile(`^-?\d{1,3}\.\d{6},\s*-?\d{1,3}\.\d{6}$`)
	return coordRegex.MatchString(input)
}

func handleCoords(input string) *models.Location {
	var (
		loc = &models.Location{}
	)
	coords := strings.Split(input, ",")
	loc.Lat = strings.TrimSpace(coords[0])
	loc.Lng = strings.TrimSpace(coords[1])
	return loc
}

func getUserLocs(chatID int64) ([]models.Location, error) {
	var locs []models.Location
	err := models.GetDB().Where("user_id = ?", chatID).Find(&locs).Error
	if err != nil {
		slog.Error("get all user locs from db", "error", err)
		return nil, err
	}
	return locs, nil
}

func getLocByID(id string) (*models.Location, error) {
	var loc *models.Location
	err := models.GetDB().Where("id = ?", id).Find(loc).Error
	if err != nil {
		slog.Error("get loc by id", "error", err)
		return nil, err
	}
	// fmt.Printf("getLocByID:\n%+v\n", loc)
	return loc, nil
}

func getID(s string) string {
	ss := strings.Split(s, ":")
	if len(ss) == 2 {
		return ss[1]
	}
	slog.Error("func getID parameter", "error", fmt.Sprintf("func getID recieved invalid parameter: %s", s))
	return ""
}
