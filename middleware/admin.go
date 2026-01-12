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
		if update.CallbackQuery != nil {
			return nil
		}

		if update.Message.Text == "/admin" {
			ctx.Bot().SendMessage(ctx, tu.Message(
				tu.ID(userid),
				"Вы не являетесь администратором бота.",
			))

			return nil
		}

		return nil
	}

	user, err := storage.GetUser(userid)
	if err != nil {
		return err
	}

	switch {
	case user.State == "nothing":
		return ctx.Next(update)
	case user.State == "awaiting_create_promocode":
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

			keyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("managePromocodes"),
				),
			)

			msg := tu.Message(
				tu.ID(userid),
				fmt.Sprintf("🎟 Вы успешно создали промокод *%s*, на *%v₽*, на %v использований.\n\nИстечёт: *%v*", code, reward/100, maxUses, expiresAt.Format("02 Jan 2006 15:04")),
			).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

			ctx.Bot().SendMessage(ctx, msg)
		}

	case user.State == "awaiting_create_category":
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "manageCategories" {
				storage.SetUserState(userid, "nothing")
				return ctx.Next(update)
			}

			return nil
		}

		if update.Message.Text != "" && update.Message != nil {
			err := storage.AddCategory(update.Message.Text)
			if err != nil {
				return err
			}

			storage.SetUserState(userid, "nothing")

			keyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCategories"),
				),
			)

			msg := tu.Message(
				tu.ID(userid),
				fmt.Sprintf("🎟 Вы успешно создали категорию *%s*", update.Message.Text),
			).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

			ctx.Bot().SendMessage(ctx, msg)
		}
	case strings.HasPrefix(user.State, "awaiting_new_product"):
		if update.CallbackQuery != nil {
			if strings.Contains(update.CallbackQuery.Data, "productsCategoryManage") {
				storage.SetUserState(userid, "nothing")
				return ctx.Next(update)
			}

			return nil
		}

		if update.Message.Text != "" && update.Message != nil {
			parts := strings.Split(user.State, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid user state format")
			}

			categoryID, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			data := strings.Split(update.Message.Text, "|")
			if len(data) < 3 {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}

			price, err := strconv.Atoi(data[2])
			if err != nil {
				ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(userid), "неверный формат"))
				return err
			}

			err = storage.AddProduct(categoryID, data[0], data[1], int64(price))
			if err != nil {
				return err
			}

			storage.SetUserState(userid, "nothing")

			keyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productsCategoryManage:" + strconv.Itoa(categoryID)),
				),
			)

			msg := tu.Message(
				tu.ID(userid),
				fmt.Sprintf("🎟 Вы успешно создали товар *%s*", update.Message.Text),
			).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

			ctx.Bot().SendMessage(ctx, msg)
		}
	case strings.HasPrefix(user.State, "awaiting_new_item"):
		if update.CallbackQuery != nil {
			if strings.Contains(update.CallbackQuery.Data, "listItems") {
				storage.SetUserState(userid, "nothing")
				return ctx.Next(update)
			}

			return nil
		}

		if update.Message.Text != "" && update.Message != nil {
			parts := strings.Split(user.State, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid user state format")
			}

			productid, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			err = storage.AddItem(productid, update.Message.Text)
			if err != nil {
				return err
			}

			storage.SetUserState(userid, "nothing")

			keyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("listItems:" + parts[1]),
				),
			)

			msg := tu.Message(
				tu.ID(userid),
				fmt.Sprintf("🎟 Вы успешно создали элемент товара *%s*", update.Message.Text),
			).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

			ctx.Bot().SendMessage(ctx, msg)
		}
	}
	return nil
}
