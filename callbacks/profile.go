package callbacks

import (
	"TelegramShop/storage"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackRefreshProfile(ctx *th.Context, query telego.CallbackQuery) error {
	chatID := query.From.ID
	username := query.From.Username
	firstname := query.From.FirstName
	lastname := query.From.LastName
	lang_code := query.From.LanguageCode

	user, err := storage.RefreshUser(chatID, username, firstname, lastname, lang_code)
	if err != nil {
		return err
	}

	editMsg := tu.EditMessageText(
		tu.ID(chatID),
		query.Message.Message().MessageID,
		fmt.Sprintf("*Профиль %s:*\n\nID: %d\nЯзык: %s\nБаланс: %d₽\nРоль: %s",
			user.Firstname,
			user.ID,
			user.LangCode,
			user.Balance,
			user.Role)).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(query.Message.Message().ReplyMarkup)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Обновлено!"))
}

func CallbackPromoCode(ctx *th.Context, query telego.CallbackQuery) error {
	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Функция в разработке!"))
}
