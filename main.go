package main

import (
	"TelegramShop/handlers"
	"TelegramShop/storage"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error load .env file:", err)
	}

	err = storage.OpenSQLite()
	if err != nil {
		log.Fatal("error to open sqlite:", err)
	}

	ctx := context.Background()
	botToken := os.Getenv("TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	bh, _ := th.NewBotHandler(bot, updates)
	defer func() { _ = bh.Stop() }()

	bh.Handle(handlers.SendMainMenu, th.CommandEqual("start"))

	bh.Handle(handlers.SendCatalog, th.TextEqual("🛍 Каталог"))
	bh.Handle(handlers.SendProfile, th.TextEqual("👤 Профиль"))
	bh.Handle(handlers.SendDeposit, th.TextEqual("💳 Пополнить баланс"))
	bh.Handle(handlers.SendSupport, th.TextEqual("🆘 Поддержка"))

	bh.HandleCallbackQuery(handlers.CallbackNextPageCat, th.CallbackDataContains("nextPageCat"))
	bh.HandleCallbackQuery(handlers.CallbackPrevPageCat, th.CallbackDataContains("prevPageCat"))
	bh.HandleCallbackQuery(handlers.CallbackNextPage, th.CallbackDataContains("nextPage"))
	bh.HandleCallbackQuery(handlers.CallbackPrevPage, th.CallbackDataContains("prevPage"))
	bh.HandleCallbackQuery(handlers.CallbackCancelCat, th.CallbackDataEqual("cancelCat"))
	bh.HandleCallbackQuery(handlers.CallbackCancel, th.CallbackDataEqual("cancel"))
	bh.HandleCallbackQuery(handlers.CallbackCategory, th.CallbackDataContains("category"))
	bh.HandleCallbackQuery(handlers.CallbackRefreshProfile, th.CallbackDataEqual("profileRefresh"))
	bh.HandleCallbackQuery(handlers.CallbackProduct, th.CallbackDataContains("product"))
	bh.HandleCallbackQuery(handlers.CallbackBuyProduct, th.CallbackDataContains("buyProduct"))
	bh.HandleCallbackQuery(handlers.CallbackBuy, th.CallbackDataContains("attentionBuy"))

	_ = bh.Start()
}
