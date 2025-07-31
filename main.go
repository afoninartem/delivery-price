package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/afoninartem/delivery-price/bot"
	"github.com/afoninartem/delivery-price/l"
)

func main() {
	l.InitLogger()
	bot.Bot()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("http listen and serve", "error", err)
		os.Exit(2)
	}
}
