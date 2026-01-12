package callbacks

import (
	"TelegramShop/storage"
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackPurchasesHistory(ctx *th.Context, query telego.CallbackQuery) error {
	page := 1

	pages, err := storage.GetPagesForPurchasesHistory(query.From.ID)
	if err != nil {
		return err
	}

	history, err := storage.GetPurchasesHistory(query.From.ID)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("У вас нет купленных товаров").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, purchase := range history {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(purchase.Name).WithCallbackData(fmt.Sprintf("purchase:%d", purchase.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePurchases:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePurchases:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите товар в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackPrevPagePurchases(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")

	pages, err := strconv.Atoi(data[2])
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	if page < 1 {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Несуществующая страница"))
	}

	history, err := storage.GetPurchasesHistory(query.From.ID)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, purchase := range history {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(purchase.Name).WithCallbackData(fmt.Sprintf("purchase:%d", purchase.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePurchases:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePurchases:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Ваши купленные товары:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackNextPagePurchases(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")

	pages, err := strconv.Atoi(data[2])
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	if page > pages {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Несуществующая страница"))
	}

	history, err := storage.GetPurchasesHistory(query.From.ID)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, purchase := range history {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(purchase.Name).WithCallbackData(fmt.Sprintf("purchase:%d", purchase.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPagePurchases:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPagePurchases:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Ваши купленные товары:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackPurchase(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	purchaseid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	purchase, err := storage.GetPurchase(purchaseid)
	if err != nil {
		return err
	}
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("purchasesHistory"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Вы купили товар №%d\n\n*Название:* %s\n*Описание:* %s\n\nЦена: *%v*₽\n\nДанные о покупке: %s", purchase.ItemID, purchase.Name, purchase.Description, purchase.Price/100, purchase.Data),
	).WithReplyMarkup(keyboard).WithParseMode(telego.ModeMarkdown)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
