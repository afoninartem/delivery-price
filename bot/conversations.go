package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/afoninartem/delivery-price/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	basePriceUrl = os.Getenv("basePriceUrl")
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
		fmt.Println("rename hitted")
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
	var loc = &models.Location{}
	err := models.GetDB().Where("id = ?", id).First(loc).Error
	if err != nil {
		fmt.Println(err)
		slog.Error("get loc by id", "error", err)
		return nil, err
	}
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

func getPrices(chatID int64, locs []models.Location) []models.Location {
	var prices = make([]models.Location, len(locs))
	for i, loc := range locs {
		prms := fmt.Sprintf("%s,%s~%s,%s", loc.Lng, loc.Lat, loc.Lng, loc.Lat)
		url := basePriceUrl + prms
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			slog.Error("new request", "error", err)
			continue
		}
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("get response", "error", err)
			continue
		}
		defer resp.Body.Close()

		b, _ := io.ReadAll(resp.Body)

		var rawData models.PriceResponse
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&rawData)
		if err != nil {
			slog.Error(string(b)[:100], "error", err)
			continue
		}
		price := extractPrice(rawData)
		loc.Price = price
		prices[i] = loc
	}

	for i, p := range prices {
		if ulp, e := userLastPrices[chatID]; e {

			prices[i].LastPrice = ulp[p.ID]
			//set new price as last seen price
			ulp[p.ID] = p.Price
			userLastPrices[chatID] = ulp
		} else {
			//set new price as last seen price
			userLastPrices[chatID] = make(map[uint]string)
			userLastPrices[chatID][p.ID] = p.Price
		}
	}
	return prices
}

func extractPrice(rawData models.PriceResponse) string {
	var price string
	for _, offer := range rawData.ClaimsOffers {
		if offer.TariffInfo.Vertical == "express" && offer.TariffInfo.Tariff == "sdd_long" {
			price = offer.Price.TotalPriceWithVat
			p, err := strconv.Atoi(price)
			if err != nil {
				slog.Error("conv string to int", "error", err)
				return price
			}
			price = fmt.Sprint(p + 20)
			break
		}
	}
	return price
}
