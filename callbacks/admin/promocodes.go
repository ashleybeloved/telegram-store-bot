package adminCallbacks

import (
	"TelegramShop/storage"
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

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
		"🎫 *Создание промокода*\n\nОтправьте данные через пробел:\n`КОД СУММА АКТИВАЦИИ ЧАСЫ`\n\n*Параметры:*\n1. *Код* — название (напр. `GIFT2026`)\n2. *Сумма* — бонус в копейках\n3. *Активации* — кол-во штук\n4. *Часы* — время жизни\n\n*Пример:*\n`PROMO100 10000 50 12`",
	).WithReplyMarkup(keyboard).WithParseMode(telego.ModeMarkdown)

	err := storage.SetUserState(query.From.ID, "awaiting_create_promocode")
	if err != nil {
		return err
	}

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackAllpromocodes(ctx *th.Context, query telego.CallbackQuery) error {
	page := 1

	pages, err := storage.GetPagesForPromocodes()
	if err != nil {
		return err
	}

	promocodes, err := storage.GetPromocodes(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Промокодов нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, promocode := range promocodes {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(promocode.Code).WithCallbackData(fmt.Sprintf("promocodeAdmin:%d", promocode.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePromocode:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePromocode:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("managePromocodes"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите промокод:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackNextPagePromocode(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	if len(data) < 3 {
		return nil
	}
	page, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	pages, err := strconv.Atoi(data[2])
	if err != nil {
		return err
	}

	if page > pages {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Несуществующая страница"))
	}

	promocodes, err := storage.GetPromocodes(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, promocode := range promocodes {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(promocode.Code).WithCallbackData(fmt.Sprintf("promocodeAdmin:%d", promocode.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePromocode:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePromocode:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("managePromocodes"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите промокод:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackPrevPagePromocode(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	if len(data) < 3 {
		return nil
	}
	page, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	pages, err := strconv.Atoi(data[2])
	if err != nil {
		return err
	}

	if page < 1 {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Несуществующая страница"))
	}

	promocodes, err := storage.GetPromocodes(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, promocode := range promocodes {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(promocode.Code).WithCallbackData(fmt.Sprintf("promocodeAdmin:%d", promocode.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePromocode:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePromocode:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("managePromocodes"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите промокод:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackPromocodeAdmin(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	promocodeid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	promocode, err := storage.GetPromocode(promocodeid)
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🗑️ Удалить").WithCallbackData(fmt.Sprint("deletePromocode:", promocode.ID)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("allPromocodes"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("🎫 *Промокод:* %s\n\n*Бонус:* %d₽\n*Активаций:* %d\n*Осталось активаций:* %d\n*Истекает через:* %s\n*Создан:* %s", promocode.Code, promocode.Reward/100, promocode.MaxUses, promocode.UsesLeft, promocode.ExpiresAt.Format("02 Jan 2006 15:04"), promocode.CreatedAt.Format("02 Jan 2006 15:04")),
	).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackDeletePromocode(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	promocodeid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	err = storage.DeletePromocode(promocodeid)
	if err != nil {
		return err
	}

	page := 1

	pages, err := storage.GetPagesForPromocodes()
	if err != nil {
		return err
	}

	promocodes, err := storage.GetPromocodes(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, promocode := range promocodes {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(promocode.Code).WithCallbackData(fmt.Sprintf("promocodeAdmin:%d", promocode.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePromocode:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePromocode:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("managePromocodes"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите промокод:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Промокод успешно удален"))
}
