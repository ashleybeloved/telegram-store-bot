package main

import (
	adminCallbacks "TelegramShop/callbacks/admin"
	callbacks "TelegramShop/callbacks/user"
	"TelegramShop/handlers"
	"TelegramShop/middleware"
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

	// User Middleware (adds to DB if not exists and checks state)

	bh.Use(middleware.UserMiddleware)

	// Handlers

	bh.Handle(handlers.SendMainMenu, th.CommandEqual("start"))

	bh.Handle(handlers.SendCatalog, th.TextEqual("🛍 Каталог"))
	bh.Handle(handlers.SendProfile, th.TextEqual("👤 Профиль"))
	bh.Handle(handlers.SendPurchasesHistory, th.TextEqual("🛒 Мои покупки"))
	bh.Handle(handlers.SendDeposit, th.TextEqual("💳 Пополнить баланс"))
	bh.Handle(handlers.SendSupport, th.TextEqual("🆘 Поддержка"))

	// Callbacks

	bh.HandleCallbackQuery(callbacks.CallbackNextPageCat, th.CallbackDataContains("nextPageCat:"))
	bh.HandleCallbackQuery(callbacks.CallbackPrevPageCat, th.CallbackDataContains("prevPageCat:"))
	bh.HandleCallbackQuery(callbacks.CallbackNextPage, th.CallbackDataContains("nextPage:"))
	bh.HandleCallbackQuery(callbacks.CallbackPrevPage, th.CallbackDataContains("prevPage:"))
	bh.HandleCallbackQuery(callbacks.CallbackCategory, th.CallbackDataContains("category:"))
	bh.HandleCallbackQuery(callbacks.CallbackProduct, th.CallbackDataContains("product:"))
	bh.HandleCallbackQuery(callbacks.CallbackBuyProduct, th.CallbackDataContains("buyProduct:"))
	bh.HandleCallbackQuery(callbacks.CallbackBuy, th.CallbackDataContains("attentionBuy:"))
	bh.HandleCallbackQuery(callbacks.CallbackCancelCat, th.CallbackDataEqual("cancelCat"))
	bh.HandleCallbackQuery(callbacks.CallbackCancel, th.CallbackDataEqual("cancel"))

	bh.HandleCallbackQuery(callbacks.CallbackRefreshProfile, th.CallbackDataEqual("profileRefresh"))

	bh.HandleCallbackQuery(callbacks.CallbackPromoCode, th.CallbackDataEqual("promoCode"))
	bh.HandleCallbackQuery(callbacks.CallbackCancelPromocode, th.CallbackDataEqual("cancelPromocode"))

	bh.HandleCallbackQuery(callbacks.CallbackPurchasesHistory, th.CallbackDataContains("purchasesHistory"))
	bh.HandleCallbackQuery(callbacks.CallbackPurchase, th.CallbackDataContains("purchase:"))
	bh.HandleCallbackQuery(callbacks.CallbackNextPagePurchases, th.CallbackDataContains("nextPagePurchases:"))
	bh.HandleCallbackQuery(callbacks.CallbackPrevPagePurchases, th.CallbackDataContains("prevPagePurchases:"))

	// Admin Middlewarec, Callbacks and Handlers

	bh.Use(middleware.AdminMiddleware)

	bh.HandleCallbackQuery(adminCallbacks.CallbackAdminMenu, th.CallbackDataEqual("adminMenu"))

	bh.HandleCallbackQuery(adminCallbacks.CallbackManagePromocodes, th.CallbackDataEqual("managePromocodes"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackCreatePromocode, th.CallbackDataEqual("createPromocode"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackAllpromocodes, th.CallbackDataEqual("allPromocodes"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNextPagePromocode, th.CallbackDataContains("nextPagePromocode:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackPrevPagePromocode, th.CallbackDataContains("prevPagePromocode:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackPromocodeAdmin, th.CallbackDataContains("promocodeAdmin:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackDeletePromocode, th.CallbackDataContains("deletePromocode:"))

	bh.HandleCallbackQuery(adminCallbacks.CallbackManageCatalog, th.CallbackDataContains("manageCatalog"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackManageCategories, th.CallbackDataContains("manageCategories"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNextPageCat, th.CallbackDataContains("nextPageCat:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackPrevPageCat, th.CallbackDataContains("prevPageCat:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackCategoryEdit, th.CallbackDataContains("categoryEdit:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackCategoryDelete, th.CallbackDataContains("categoryDelete:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackCategoryCreate, th.CallbackDataContains("categoryCreate"))

	bh.HandleCallbackQuery(adminCallbacks.CallbackManageProducts, th.CallbackDataContains("manageProducts"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNewProduct, th.CallbackDataContains("newProduct:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNextPage, th.CallbackDataContains("nextPage:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackPrevPage, th.CallbackDataContains("prevPage:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackCategory, th.CallbackDataContains("productsCategoryManage:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackProductManage, th.CallbackDataContains("productManage:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackDeleteProduct, th.CallbackDataContains("deleteProduct:"))

	bh.HandleCallbackQuery(adminCallbacks.CallbackListItems, th.CallbackDataContains("listItems:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackItemManage, th.CallbackDataContains("itemManage:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNewItem, th.CallbackDataContains("newItem:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackItemDelete, th.CallbackDataContains("itemDelete:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackNextPageItems, th.CallbackDataContains("nextPageItems:"))
	bh.HandleCallbackQuery(adminCallbacks.CallbackPrevPageItems, th.CallbackDataContains("prevPageItems:"))

	bh.Handle(handlers.SendAdminMenu, th.CommandEqual("admin"))

	log.Println("Bot started")

	_ = bh.Start()
}
