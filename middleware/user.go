package middleware

import (
	"TelegramShop/storage"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func UserMiddleware(ctx *th.Context, update telego.Update) error {
	var userid int64
	var username, firstname, lastname, langCode string

	if update.CallbackQuery != nil {
		userid = update.CallbackQuery.From.ID
		username = update.CallbackQuery.From.Username
		firstname = update.CallbackQuery.From.FirstName
		lastname = update.CallbackQuery.From.LastName
		langCode = update.CallbackQuery.From.LanguageCode
	} else if update.Message != nil {
		userid = update.Message.From.ID
		username = update.Message.From.Username
		firstname = update.Message.From.FirstName
		lastname = update.Message.From.LastName
		langCode = update.Message.From.LanguageCode
	} else {
		return ctx.Next(update)
	}

	err := storage.AddUser(userid, username, firstname, lastname, langCode)
	if err != nil {
		return err
	}

	user, err := storage.GetUser(userid)
	if err != nil {
		return err
	}

	switch user.State {
	case "nothing":
		return ctx.Next(update)
	case "awaiting_promocode":
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "cancelPromocode" {
				storage.SetUserState(userid, "nothing")
				return ctx.Next(update)
			}

			return nil
		}

		if update.Message != nil && update.Message.Text != "" {
			reward, err := storage.RedeemPromocode(userid, update.Message.Text)
			if err != nil {
				keyboard := tu.InlineKeyboard(
					tu.InlineKeyboardRow(
						tu.InlineKeyboardButton("⬅️ Отмена").WithCallbackData("cancelPromocode"),
					),
				)

				ctx.Bot().SendMessage(ctx, tu.Message(
					tu.ID(userid),
					fmt.Sprintf("❌ Ошибка при активации промокода: %v", err),
				).WithReplyMarkup(keyboard))
				return nil
			}

			storage.SetUserState(userid, "nothing")

			msg := tu.Message(
				tu.ID(update.Message.From.ID),
				fmt.Sprintf("Вы успешно активировали промокод %s, на %v₽", update.Message.Text, reward),
			)

			ctx.Bot().SendMessage(ctx, msg)
		}
	default:
		return ctx.Next(update)
	}

	return nil
}
