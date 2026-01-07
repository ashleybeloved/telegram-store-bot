package handlers

import (
	"TelegramShop/storage"
	"strconv"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMainMenu(ctx *th.Context, update telego.Update) (err error) {
	chatID := update.Message.Chat.ID
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
			tu.KeyboardButton("🛒 Корзина"),
			tu.KeyboardButton("👤 Профиль"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("💳 Пополнить баланс"),
			tu.KeyboardButton("🆘 Поддержка"),
		),
	).WithResizeKeyboard()

	msg := tu.Photo(
		tu.ID(chatID),
		tu.FileFromID(photo),
	).WithCaption(firstname + ", добро пожаловать в <b>heaven.help</b>").WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard)

	ctx.Bot().SendPhoto(ctx, msg)

	return nil
}

func SendCatalog(ctx *th.Context, update telego.Update) (err error) {
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Не реализовано =(",
	)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendCart(ctx *th.Context, update telego.Update) (err error) {
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Не реализовано =(",
	)

	ctx.Bot().SendMessage(ctx, msg)

	return nil
}

func SendProfile(ctx *th.Context, update telego.Update) (err error) {
	user, err := storage.FindUser(update.Message.From.ID)
	if err != nil {
		return err
	}

	chatID := update.Message.Chat.ID

	msg := tu.Message(
		tu.ID(chatID),
		"<b>Профиль "+user.Firstname+":</b>\n\nID: "+strconv.Itoa(user.ID)+"\nЯзык: "+user.LangCode+"\nБаланс: "+strconv.FormatInt(user.Balance, 10)+"₽"+"\nРоль: "+user.Role).WithParseMode(telego.ModeHTML)

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
