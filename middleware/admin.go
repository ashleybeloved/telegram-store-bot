package middleware

import (
	"TelegramShop/storage"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func AdminMiddleware(ctx *th.Context, update telego.Update) error {
	var userid int64
	adminid := os.Getenv("ADMIN_ID")

	if update.CallbackQuery != nil {
		userid = update.CallbackQuery.From.ID
	} else if update.Message != nil {
		userid = update.Message.From.ID
	} else {
		return nil
	}

	if fmt.Sprint(userid) != adminid {
		ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(userid),
			"Вы не являетесь администратором бота.",
		))
		return nil
	}

	user, err := storage.GetUser(userid)
	if err != nil {
		return err
	}

	switch user.State {
	case "nothing":
		return ctx.Next(update)
	case "awaiting_create_promocode":
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "managePromocodes" {
				storage.SetUserState(userid, "nothing")
				return ctx.Next(update)
			}

			return nil
		}

		if update.Message.Text != "" && update.Message != nil {
			data := strings.Split(update.Message.Text, " ")
			if len(data) < 4 {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}
			code := data[0]
			if code == "" {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}
			reward, err := strconv.ParseInt(data[1], 10, 64)
			if err != nil {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}
			maxUses, err := strconv.Atoi(data[2])
			if err != nil {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}
			hoursCount, err := strconv.Atoi(data[3])
			if err != nil {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}

			expiresAt := time.Now().Add(time.Duration(hoursCount) * time.Hour)

			err = storage.NewPromocode(code, reward, maxUses, expiresAt)
			if err != nil {
				return err
			}

			storage.SetUserState(userid, "nothing")

			msg := tu.Message(
				tu.ID(userid),
				fmt.Sprintf("Вы успешно создали промокод *%s*, на *%v₽*, на %v использований.\n\nИстечёт: *%v*", code, reward, maxUses, expiresAt),
			).WithParseMode(telego.ModeMarkdown)

			ctx.Bot().SendMessage(ctx, msg)
		}
	}

	return ctx.Next(update)
}
