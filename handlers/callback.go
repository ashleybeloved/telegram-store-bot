package handlers

import (
	"TelegramShop/storage"
	"fmt"
	"strconv"
	"strings"

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
		fmt.Sprintf("<b>Профиль %s:</b>\n\nID: %d\nЯзык: %s\nБаланс: %d₽\nРоль: %s",
			user.Firstname,
			user.ID,
			user.LangCode,
			user.Balance,
			user.Role)).WithParseMode(telego.ModeHTML).WithReplyMarkup(query.Message.Message().ReplyMarkup)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Обновлено!"))
}

func CallbackPrevPage(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	pageStr := data[1]
	pagesStr := data[2]
	pages, _ := strconv.Atoi(pagesStr)
	page, _ := strconv.Atoi(pageStr)

	if page < 1 {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Не существующая страница"))
	}

	pages, err := storage.GetPagesForCategories()
	if err != nil {
		return err
	}

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("cat:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите категорию товаров:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackNextPage(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	pageStr := data[1]
	pagesStr := data[2]
	pages, _ := strconv.Atoi(pagesStr)
	page, _ := strconv.Atoi(pageStr)

	if page > pages {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Не существующая страница"))
	}

	pages, err := storage.GetPagesForCategories()
	if err != nil {
		return err
	}

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("cat:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancel"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите категорию товаров:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
