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

func CallbackManageCatalog(ctx *th.Context, query telego.CallbackQuery) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📃 Категории").WithCallbackData("manageCategories"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🛒 Товары").WithCallbackData("manageProducts"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("adminMenu"),
		))

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Редактирование каталога:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackManageCategories(ctx *th.Context, query telego.CallbackQuery) error {
	page := 1

	pages, err := storage.GetPagesForCategories()
	if err != nil {
		return err
	}

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Категорий нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("categoryEdit:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать категорию").WithCallbackData("categoryCreate"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCatalog"),
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

func CallbackPrevPageCat(ctx *th.Context, query telego.CallbackQuery) error {
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

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("categoryEdit:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать категорию").WithCallbackData("categoryCreate"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCategories"),
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

func CallbackNextPageCat(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	pages, _ := strconv.Atoi(data[2])
	page, _ := strconv.Atoi(data[1])

	if page > pages {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Несуществующая страница"))
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
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("categoryEdit:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать категорию").WithCallbackData("categoryCreate"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCategories"),
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

func CallbackCategoryEdit(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	cat_id, _ := strconv.Atoi(data[1])

	category, err := storage.GetCategory(cat_id)
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🗑️ Удалить категорию").WithCallbackData(fmt.Sprintf("categoryDelete:%d", category.ID)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCategories"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Категория ID: %d\nИмя: %v\nПродукты: %v", category.ID, category.Name, category.Products),
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackCategoryCreate(ctx *th.Context, query telego.CallbackQuery) error {
	err := storage.SetUserState(query.From.ID, "awaiting_create_category")
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCategories"),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Введите название категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackCategoryDelete(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	catid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	err = storage.DelCategory(catid)
	if err != nil {
		return err
	}

	page := 1

	pages, err := storage.GetPagesForCategories()
	if err != nil {
		return err
	}

	categories, err := storage.GetCategories(page)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Категорий нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, cat := range categories {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("categoryEdit:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать категорию").WithCallbackData("categoryCreate"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageCatalog"),
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
