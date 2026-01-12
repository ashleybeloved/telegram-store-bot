package handlers

import (
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendAdminMenu(ctx *th.Context, update telego.Update) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📦 Управление каталогом").WithCallbackData("manageCatalog"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🎟 Управление промокодами").WithCallbackData("managePromocodes"),
		))

	ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Админ-панель",
	).WithReplyMarkup(keyboard))

	return nil
}
