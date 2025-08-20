package main

import (
	"github.com/afoninartem/delivery-price/bot"
	"github.com/afoninartem/delivery-price/l"
)

func main() {
	l.InitLogger()
	bot.Bot()
}
