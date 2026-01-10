package callbacks

import (
	"TelegramShop/storage"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackRefreshProfile(ctx *th.Context, query telego.CallbackQuery) error {
	user, err := storage.RefreshUser(query.From.ID, query.From.Username, query.From.FirstName, query.From.LastName, query.From.LanguageCode)
	if err != nil {
		return err
	}

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
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
