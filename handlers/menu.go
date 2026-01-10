package handlers

import (
	"TelegramShop/configs"
	"TelegramShop/storage"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMainMenu(ctx *th.Context, update telego.Update) (err error) {
	photo := configs.MainMenuPhotoID

	err = storage.SetUserState(update.Message.From.ID, "nothing")
	if err != nil {
		return err
	}

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("🛍 Каталог"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🛒 Мои покупки"),
			tu.KeyboardButton("👤 Профиль"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("💳 Пополнить баланс"),
			tu.KeyboardButton("🆘 Поддержка"),
		),
	).WithResizeKeyboard()

	msg := tu.Photo(
		tu.ID(update.Message.From.ID),
		tu.FileFromID(photo),
	).WithCaption(update.Message.From.FirstName + ", добро пожаловать в *heaven.help*").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().SendPhoto(ctx, msg)

	return nil
}

func SendCatalog(ctx *th.Context, update telego.Update) (err error) {
	page := 1

	pages, err := storage.GetPagesForCategories()
	if err != nil {
		return err
	}

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("category:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Выберите категорию товаров:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendProfile(ctx *th.Context, update telego.Update) (err error) {
	user, err := storage.GetUser(update.Message.From.ID)
	if err != nil {
		return err
	}

	chatID := update.Message.Chat.ID

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔄 Обновить").WithCallbackData("profileRefresh"),
		),
	)

	msg := tu.Message(
		tu.ID(chatID),
		fmt.Sprintf("**Профиль %s:**\n\nID: %d\nЯзык: %s\nБаланс: %d₽\nРоль: %s",
			user.Firstname,
			user.ID,
			user.LangCode,
			user.Balance,
			user.Role)).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendDeposit(ctx *th.Context, update telego.Update) (err error) {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🎟 Промокод").WithCallbackData("promoCode"),
		),
	)

	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Пополнение баланса доступно через методы ниже:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendSupport(ctx *th.Context, update telego.Update) (err error) {
	chatID := update.Message.Chat.ID
	photo := configs.SupportPhotoID

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Support").WithURL("t.me/fuckcensor"),
		),
	)

	msg := tu.Photo(
		tu.ID(chatID),
		tu.FileFromID(photo),
	).WithCaption("<b>Нужна помощь?</b>\n\nМожете обратиться в саппорт!").WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard)

	ctx.Bot().SendPhoto(ctx, msg)

	return nil
}

func SendPurchases(ctx *th.Context, update telego.Update) (err error) {

	return nil
}
