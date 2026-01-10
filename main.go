package main

import (
	"TelegramShop/callbacks"
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

	// Handlers

	bh.Handle(handlers.SendMainMenu, th.CommandEqual("start"))

	bh.Handle(handlers.SendCatalog, th.TextEqual("🛍 Каталог"))
	bh.Handle(handlers.SendProfile, th.TextEqual("👤 Профиль"))
	bh.Handle(handlers.SendPurchases, th.TextEqual("🛒 Мои покупки"))
	bh.Handle(handlers.SendDeposit, th.TextEqual("💳 Пополнить баланс"))
	bh.Handle(handlers.SendSupport, th.TextEqual("🆘 Поддержка"))

	// Callbacks

	bh.HandleCallbackQuery(callbacks.CallbackNextPageCat, th.CallbackDataContains("nextPageCat"))
	bh.HandleCallbackQuery(callbacks.CallbackPrevPageCat, th.CallbackDataContains("prevPageCat"))
	bh.HandleCallbackQuery(callbacks.CallbackNextPage, th.CallbackDataContains("nextPage"))
	bh.HandleCallbackQuery(callbacks.CallbackPrevPage, th.CallbackDataContains("prevPage"))
	bh.HandleCallbackQuery(callbacks.CallbackCategory, th.CallbackDataContains("category"))
	bh.HandleCallbackQuery(callbacks.CallbackProduct, th.CallbackDataContains("product"))
	bh.HandleCallbackQuery(callbacks.CallbackBuyProduct, th.CallbackDataContains("buyProduct"))
	bh.HandleCallbackQuery(callbacks.CallbackBuy, th.CallbackDataContains("attentionBuy"))
	bh.HandleCallbackQuery(callbacks.CallbackCancelCat, th.CallbackDataEqual("cancelCat"))
	bh.HandleCallbackQuery(callbacks.CallbackCancel, th.CallbackDataEqual("cancel"))

	bh.HandleCallbackQuery(callbacks.CallbackRefreshProfile, th.CallbackDataEqual("profileRefresh"))
	bh.HandleCallbackQuery(callbacks.CallbackPromoCode, th.CallbackDataEqual("promoCode"))

	log.Println("Bot started")

	_ = bh.Start()
}
