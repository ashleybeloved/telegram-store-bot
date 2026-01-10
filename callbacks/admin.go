package callbacks

import (
	"TelegramShop/storage"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackAdminMenu(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📦 Управление товарами").WithCallbackData("manageProducts"),
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

func CallbackManagePromocodes(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать промокод").WithCallbackData("createPromocode"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🧾 Список промокодов").WithCallbackData("allPromocodes"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("adminMenu"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Управление промокодами:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackCreatePromocode(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Отменить").WithCallbackData("managePromocodes"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"🎫 *Создание промокода*\n\nОтправьте данные через пробел:\n`КОД СУММА АКТИВАЦИИ ЧАСЫ`\n\n*Параметры:*\n1. *Код* — название (напр. `GIFT2026`)\n2. *Сумма* — бонус в рублях\n3. *Активации* — кол-во штук\n4. *Часы* — время жизни\n\n*Пример:*\n`PROMO100 100 50 12`",
	).WithReplyMarkup(keyboard).WithParseMode(telego.ModeMarkdown)

	err := storage.SetUserState(query.From.ID, "awaiting_create_promocode")
	if err != nil {
		return err
	}

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
