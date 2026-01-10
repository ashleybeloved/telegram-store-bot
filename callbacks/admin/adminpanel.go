package adminCallbacks

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackAdminMenu(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📦 Управление каталогом").WithCallbackData("manageCatalog"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🎟 Управление промокодами").WithCallbackData("managePromocodes"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("👥 Управление пользователями").WithCallbackData("manageUsers"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📊 Статистика").WithCallbackData("viewStats"),
		))

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Админ-панель",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
