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

func CallbackManageProducts(ctx *th.Context, query telego.CallbackQuery) error {
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
			tu.InlineKeyboardButton(cat.Name).WithCallbackData(fmt.Sprintf("productsCategoryManage:%d", cat.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageCat:%d:%d", page-1, pages)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageCat:%d:%d", page+1, pages)),
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
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("productManage:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать товар").WithCallbackData("newProduct:"+strconv.Itoa(cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageProducts"),
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
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("productManage:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать товар").WithCallbackData("newProduct:"+strconv.Itoa(cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageProducts"),
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
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("productManage:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать товар").WithCallbackData("newProduct:"+strconv.Itoa(cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageProducts"),
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

func CallbackProductManage(ctx *th.Context, query telego.CallbackQuery) error {
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
			tu.InlineKeyboardButton("➕ Пополнить").WithCallbackData("newItem:"+strconv.Itoa(int(product.ID))),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📋 Список товаров").WithCallbackData("listItems:"+strconv.Itoa(int(product.ID))),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🗑️ Удалить").WithCallbackData("deleteProduct:"+strconv.Itoa(int(product.ID))),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productsCategoryManage:"+strconv.Itoa(int(product.CategoryID))),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Товар ID: %d\nИмя: %v\nОписание: %v\nЦена: %v\nВ наличии: %d шт.", product.ID, product.Name, product.Description, product.Price, product.Stock),
	).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackNewProduct(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	cat_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	err = storage.SetUserState(query.From.ID, fmt.Sprintf("awaiting_new_product:%d", cat_id))
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productsCategoryManage:" + strconv.Itoa(cat_id)),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Для создания товара в базе, отправьте боту сообщение строго в следующем формате:\n\nФОРМАТ: *Название|Описание|Цена(!В КОПЕЙКАХ!)*\n\nПРИМЕР: Подписка Telegram Premium|1 месяц, активация подарком|45000",
	).WithReplyMarkup(keyboard).WithParseMode(telego.ModeMarkdown)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackDeleteProduct(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	product_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	product, err := storage.GetProduct(product_id)
	if err != nil {
		return err
	}

	cat_id := product.CategoryID

	err = storage.DelProduct(product_id)
	if err != nil {
		return err
	}

	page := 1

	pages, err := storage.GetPagesForProducts(int(cat_id))
	if err != nil {
		return err
	}

	products, err := storage.GetProducts(page, int(cat_id))
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
			tu.InlineKeyboardButton(product.Name).WithCallbackData(fmt.Sprintf("productManage:%d", product.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPage:%d:%d:%d", page-1, pages, cat_id)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPage:%d:%d:%d", page+1, pages, cat_id)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать товар").WithCallbackData("newProduct:"+strconv.Itoa(int(cat_id))),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("manageProducts"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите товар в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Товар успешно удалён"))
}

func CallbackListItems(ctx *th.Context, query telego.CallbackQuery) error {
	page := 1
	data := strings.Split(query.Data, ":")
	productid, _ := strconv.Atoi(data[1])

	pages, err := storage.GetPagesForItems(productid)
	if err != nil {
		return err
	}

	items, err := storage.GetItems(page, productid)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Элементов нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, item := range items {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprint(item.ID)).WithCallbackData(fmt.Sprintf("itemManage:%d", item.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageItems:%d:%d:%d", page-1, pages, productid)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageItems:%d:%d:%d", page+1, pages, productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать элемент товара").WithCallbackData("newItem:"+strconv.Itoa(productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productManage:"+strconv.Itoa(productid)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите элемент товара в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackPrevPageItems(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	productid, err := strconv.Atoi(data[3])
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

	items, err := storage.GetItems(page, productid)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, item := range items {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprint(item.ID)).WithCallbackData(fmt.Sprintf("itemManage:%d", item.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageItems:%d:%d:%d", page-1, pages, productid)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageItems:%d:%d:%d", page+1, pages, productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать элемент товара").WithCallbackData("newItem:"+strconv.Itoa(productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productManage:"+strconv.Itoa(productid)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите элемент товара в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackNextPageItems(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	productid, err := strconv.Atoi(data[3])
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

	items, err := storage.GetItems(page, productid)
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	for _, item := range items {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprint(item.ID)).WithCallbackData(fmt.Sprintf("itemManage:%d", item.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageItems:%d:%d:%d", page-1, pages, productid)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageItems:%d:%d:%d", page+1, pages, productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать элемент товара").WithCallbackData("newItem:"+strconv.Itoa(productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productManage:"+strconv.Itoa(productid)),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите элемент товара в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackItemManage(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	itemid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	item, err := storage.GetItem(itemid)
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🗑️ Удалить").WithCallbackData(fmt.Sprintf("itemDelete:%d", item.ID)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData(fmt.Sprintf("productManage:%d", item.ProductID)),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		fmt.Sprintf("Элемент товара %d\nДанные при выдаче товара: %v", item.ID, item.Data),
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
}

func CallbackItemDelete(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")

	itemid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	item, err := storage.GetItem(itemid)
	if err != nil {
		return err
	}

	productid := item.ProductID

	err = storage.DelItem(itemid)
	if err != nil {
		return err
	}

	page := 1

	pages, err := storage.GetPagesForItems(int(productid))
	if err != nil {
		return err
	}

	items, err := storage.GetItems(page, int(productid))
	if err != nil {
		return err
	}

	var rows [][]telego.InlineKeyboardButton

	if pages == 0 {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Элементов нет").WithCallbackData(" "),
		))
		pages = 1
	}

	for _, item := range items {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprint(item.ID)).WithCallbackData(fmt.Sprintf("itemManage:%d", item.ID)),
		))
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("<").WithCallbackData(fmt.Sprintf("prevPageItems:%d:%d:%d", page-1, pages, productid)),
			tu.InlineKeyboardButton(fmt.Sprintf("%d/%d", page, pages)).WithCallbackData(" "),
			tu.InlineKeyboardButton(">").WithCallbackData(fmt.Sprintf("nextPageItems:%d:%d:%d", page+1, pages, productid)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Создать элемент товара").WithCallbackData("newItem:"+strconv.Itoa(int(productid))),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("productManage:"+strconv.Itoa(int(productid))),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Выберите элемент товара в категории:",
	).WithReplyMarkup(keyboard)

	ctx.Bot().EditMessageText(ctx, editMsg)

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Элемент товара удален"))
}

func CallbackNewItem(ctx *th.Context, query telego.CallbackQuery) error {
	data := strings.Split(query.Data, ":")
	productid, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	err = storage.SetUserState(query.From.ID, fmt.Sprintf("awaiting_new_item:%d", productid))
	if err != nil {
		return err
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData("listItems:" + strconv.Itoa(productid)),
		),
	)

	editMsg := tu.EditMessageText(
		tu.ID(query.From.ID),
		query.Message.Message().MessageID,
		"Для создания элемента товара в базе, отправьте боту сообщение строго в следующем формате:\n\nФОРМАТ: *Данные элемента товара*\n\nПРИМЕР: serialkey-12345-67890",
	).WithReplyMarkup(keyboard).WithParseMode(telego.ModeMarkdown)

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
		fmt.Sprintf("Категория ID: %d\nИмя: %v", category.ID, category.Name),
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
