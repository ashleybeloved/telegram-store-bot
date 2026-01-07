package handlers

import (
	"TelegramShop/storage"
	"strconv"

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
		"<b>Профиль "+user.Firstname+":</b>\n\nID: "+strconv.Itoa(user.ID)+"\nЯзык: "+user.LangCode+"\nБаланс: "+strconv.FormatInt(user.Balance, 10)+"₽"+"\nРоль: "+user.Role).WithParseMode(telego.ModeHTML).WithReplyMarkup(query.Message.Message().ReplyMarkup)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Обновлено!"))
}
