package middleware

import (
	"TelegramShop/storage"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func UserMiddleware(ctx *th.Context, update telego.Update) error {
	storage.AddUser(update.Message.From.ID, update.Message.From.Username, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.LanguageCode)

	user, err := storage.FindUser(update.Message.From.ID)
	if err != nil {
		return err
	}

	switch user.State {
	case "nothing":
		return ctx.Next(update)
	}

	return nil
}
