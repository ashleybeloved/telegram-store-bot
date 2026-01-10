package handlers

import (
	"TelegramShop/storage"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMainMenu(ctx *th.Context, update telego.Update) (err error) {
	user_id := update.Message.From.ID
	username := update.Message.From.Username
	firstname := update.Message.From.FirstName
	lastname := update.Message.From.LastName
	lang_code := update.Message.From.LanguageCode

	storage.AddUser(user_id, username, firstname, lastname, lang_code)

	photo := "AgACAgIAAxkBAAPGaV6tpwnR1_akAyzb6MH26kzBpNgAAkgTaxuV-fBKuPW7m2HJYfIBAAMCAAN5AAM4BA"

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
		tu.ID(user_id),
		tu.FileFromID(photo),
	).WithCaption(firstname + ", добро пожаловать в *heaven.help*").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

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
	user, err := storage.FindUser(update.Message.From.ID)
	if err != nil {
		return err
	}

	chatID := update.Message.Chat.ID

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔄 Обновить").WithCallbackData("profileRefresh"),
			tu.InlineKeyboardButton("🎁 Ввести промокод").WithCallbackData("promoCode"),
		),
	)

	msg := tu.Message(
		tu.ID(chatID),
		fmt.Sprintf("<b>Профиль %s:</b>\n\nID: %d\nЯзык: %s\nБаланс: %d₽\nРоль: %s",
			user.Firstname,
			user.ID,
			user.LangCode,
			user.Balance,
			user.Role)).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendDeposit(ctx *th.Context, update telego.Update) (err error) {
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Не реализовано =(",
	)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendSupport(ctx *th.Context, update telego.Update) (err error) {
	chatID := update.Message.Chat.ID
	photo := "AgACAgIAAxkBAAPGaV6tpwnR1_akAyzb6MH26kzBpNgAAkgTaxuV-fBKuPW7m2HJYfIBAAMCAAN5AAM4BA"

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
