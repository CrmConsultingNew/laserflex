package laserflex

func LaserProductsAddToDeal() {
	// Предварительно добавляем остатки по товарам на складах.
	// Как нам получить ID товара? и добавить его в задачу. (Вероятно новым вебхуком забрать все товары с их ID)
	// Тут добавляем товары в задачу. Название и Цена из КП?

}

func LaserAddWarehouseDocument() {
	// нужно создать документ S на оприходование

	// создать документ - catalog.document.add . Как добавить в вебхук?
	// получение полей складского учета - catalog.document.getFields
	// возвращает типы - catalog.enum.getStoreDocumentTypes
	//(A – Приход товара на склад;
	//S – Оприходование товара;
	//M – Перемещение товара между складами;
	//R – Возврат товара;
	//D – Списание товара.)
}

func LaserTasksAdd(entityId, smartProcessId, elementId string) {

	// получить все пользовательские поля в задаче (названия и свойства) - tasks.task.getFields. Ищем примерно это - UF_AUTO_584432783202. Добавляем в fields.
	// Тут отправляем вебхук - https://bitrix.laser-flex.ru/rest/149/64ccuxau3kqtkfls/task.item.userfield.getfields.json
	// и забираем ответ (значения полей пользовательских)

	//task.item.add - Привязка к смарту?
	//tasks.task.add - Привязка к сделке?

	// Создаем основные задачи + Подзадачи

	// ID группы 1 - Лазерные работы
	// ID группы 11 - Труборез
	// ID группы 10 - Гибочные работы
	// ID группы 2 - Производство
	// ID группы 12 - Нанесение покрытий

	// RETURN urlCreateTask := fmt.Sprintf("https://bitrix.laser-flex.ru/rest/149/ycz7102vaerygxvb/task.item.add.json?fields[TITLE]=%s&fields[RESPONSIBLE_ID]=%s&fields[GROUP_ID]=%s]")
	// Тут должны вернуть ID созданной задачи.

	// В подзадачах указываем PARENT_ID, т.е. при создании задачи мы должны вернуть ID из битрикс. Передаем в Вебхук
	// RETURN urlCreateTaskWithParent := fmt.Sprintf("https://bitrix.laser-flex.ru/rest/149/ycz7102vaerygxvb/task.item.add.json?fields[TITLE]=%s&fields[RESPONSIBLE_ID]=%s&fields[GROUP_ID]=%s&PARENT_ID=%s]")
	// Тут тоже должны вернуть ID созданной подзадачи.

	//Привязка задачи к смарт-процессу. Читаем комменты по ссылке - https://dev.1c-bitrix.ru/rest_help/tasks/task/tasks/tasks_task_getFields.php
	// Вот функция которая преобразует в 16ричное значение - func convertEntityTypeIdToHex(entityTypeId int, elementId int) string {
	//	// Преобразуем entityTypeId в шестнадцатеричный формат
	//	hexEntityTypeId := fmt.Sprintf("%x", entityTypeId)
	//	// Формируем финальную строку с префиксом T и добавляем elementId
	//	return fmt.Sprintf("T%s_%d", hexEntityTypeId, elementId)
	//}

}
