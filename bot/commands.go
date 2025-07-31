package bot

// func commands() {
// 	cmdCfg := tgbotapi.NewSetMyCommands(
// 		tgbotapi.BotCommand{
// 			Command:     "start",
// 			Description: "Начало работы",
// 		},
// 		tgbotapi.BotCommand{
// 			Command:     "add",
// 			Description: "Добавить точку",
// 		},
// 		tgbotapi.BotCommand{
// 			Command:     "edit",
// 			Description: "Редактировать список",
// 		},
// 		tgbotapi.BotCommand{
// 			Command:     "get",
// 			Description: "Проверить цены",
// 		},
// 		tgbotapi.BotCommand{
// 			Command:     "help",
// 			Description: "Помощь",
// 		},
// 	)
// 	bot.Send(cmdCfg)
// }

// func handleCommand(cmd string, chatID int64) {
// 	var (
// 		msg = tgbotapi.NewMessage(chatID, "")
// 	)
// 	switch cmd {
// 	case "start":
// 		msg.Text = "Выберите одно из действий:"
// 		msg.ReplyMarkup = mainMenuKB()
// 	case "add":
// 		msg.Text = "Добавьте новую точку, перейдя в Mini App."
// 	case "edit":
// 		msg.Text = "Вы можете переименовать или удалить локацию:"
// 	case "get":
// 		msg.Text = "Актуальные цена на доставку в ваших локациях:"
// 	case "help":
// 		msg.Text = "Добро пожаловать в раздел помощи.\nОписание возможностей бота:\n"
// 	}
// 	bot.Send(msg)
// }
