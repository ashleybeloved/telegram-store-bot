package callbacks

import (
	"TelegramShop/configs"
	"TelegramShop/storage"
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CallbackCancelCat(ctx *th.Context, query telego.CallbackQuery) error {
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
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("category:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
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

func CallbackCancel(ctx *th.Context, query telego.CallbackQuery) error {
	photo := configs.MainMenuPhotoID

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("🛍 Каталог"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🛒 Мои покупки"),
			tu.KeyboardButton("👤 Профиль"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("💳 Пополнить баланс"),
			tu.KeyboardButton("🆘 Поддержка"),
		),
	).WithResizeKeyboard()

	msg := tu.Photo(
		tu.ID(query.From.ID),
		tu.FileFromID(photo),
	).WithCaption(query.From.FirstName + ", добро пожаловать в *heaven.help*").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().SendPhoto(ctx, msg)

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
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("category:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
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
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("category:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
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

func CallbackCategory(ctx *th.Context, query telego.CallbackQuery) error {
	page := 1
	data := strings.Split(query.Data, ":")
	cat_id, _ := strconv.Atoi(data[1])

	pages, err := storage.GetPagesForProducts(cat_id)
	if err != nil {
		return err
	}

	products, err := storage.GetProducts(page, cat_id)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Товаров нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, product := range products {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("product:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancelCat"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
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

func CallbackPrevPage(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	cat_id, err := strconv.Atoi(data[3])
	if err != nil {
		return err
	}

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

	products, err := storage.GetProducts(page, cat_id)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, product := range products {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("product:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancelCat"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
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

func CallbackNextPage(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	cat_id, err := strconv.Atoi(data[3])
	if err != nil {
		return err
	}

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

	products, err := storage.GetProducts(page, cat_id)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, product := range products {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("product:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("cancelCat"),
			tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search"),
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

func CallbackProduct(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	product_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	product, err := storage.GetProduct(product_id)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData(fmt.Sprintf("category:%d", product.CategoryID)),
			tu.InlineKeyboardButton("🛒 Купить").WithCallbackData(fmt.Sprintf("buyProduct:%d", product.ID)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Товар #%v: *%s*\n\nОписание: %s\n\nОсталось: %d\n\nЦена: *%d₽*", product.ID, product.Name, product.Description, product.Stock, product.Price/100),
	).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackBuyProduct(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	product_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	user, err := storage.GetUser(query.From.ID)
	if err != nil {
		return err
	}

	product, err := storage.GetProduct(product_id)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData(fmt.Sprintf("product:%d", product.ID)),
			tu.InlineKeyboardButton("🛒 Уверен").WithCallbackData(fmt.Sprintf("attentionBuy:%d", product.ID)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Вы уверены что хотите купить: *%s* за *%v₽*?\n\nВаш баланс: %v₽", product.Name, product.Price/100, user.Balance/100),
	).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackBuy(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	product_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	user, err := storage.GetUser(query.From.ID)
	if err != nil {
		return err
	}

	product, err := storage.GetProduct(product_id)
	if err != nil {
		return err
	}

	itemData, err := storage.BuyProduct(user.UserID, product.ID)
	if err != nil {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText(fmt.Sprint(err)))
	}

	var rows [][]telego.InlineKeyboardButton

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ В Каталог").WithCallbackData(fmt.Sprintf("cancelCat:%d", product.ID)),
			tu.InlineKeyboardButton("🛒 Мои покупки").WithCallbackData(fmt.Sprintf("category:%d", product.CategoryID)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Вы успешно приобрели: *%s*\n\nДанные приложенные к товару: %s", product.Name, itemData),
	).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}
