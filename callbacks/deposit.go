package callbacks

import (
	"TelegramShop/storage"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackPromoCode(ctx *th.Context, query telego.CallbackQuery) error {
	err := storage.SetUserState(query.From.ID, "awaiting_promocode")
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Отмена").WithCallbackData("cancelPromocode"),
		),
	)

	msg := tu.Message(
		tu.ID(query.From.ID),
		"🎁 Введите промокод для активации:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackCancelPromocode(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🎁 Промокод").WithCallbackData("promoCode"),
		),
	)

	msg := tu.Message(
		tu.ID(query.From.ID),
		"Пополнение баланса доступно через методы ниже:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().SendMessage(ctx, msg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
