package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Task struct {
	Title         string   `json:"TITLE"`
	ResponsibleID int      `json:"RESPONSIBLE_ID"`
	GroupID       int      `json:"GROUP_ID,omitempty"`
	UfCrmTask     []string `json:"UF_CRM_TASK,omitempty"`
}

// Структура для пользовательских полей задачи
type CustomTaskFields struct {
	OrderNumber       string `json:"UF_AUTO_303168834495,omitempty"` // № заказа
	Customer          string `json:"UF_AUTO_876283676967,omitempty"` // Заказчик
	Manager           string `json:"UF_AUTO_794809224848,omitempty"` // Менеджер
	Material          string `json:"UF_AUTO_468857876599,omitempty"` // Материал
	Comment           string `json:"UF_AUTO_497907774817,omitempty"` // Комментарий
	ProductionTask    string `json:"UF_AUTO_433735177517,omitempty"` // Произв. Задача
	Bend              string `json:"UF_AUTO_726724682983,omitempty"` // Гибка
	Coating           string `json:"UF_AUTO_512869473370,omitempty"` // Покрытие
	TemporaryOrderSum string `json:"UF_AUTO_555642596740,omitempty"` // Временная сумма заказа
	Quantity          string `json:"UF_AUTO_552243496167,omitempty"` // Кол-во
}

// Структура для создания задачи с полями
type TaskWithParent struct {
	Title         string           `json:"TITLE"`
	ResponsibleID int              `json:"RESPONSIBLE_ID"`
	GroupID       int              `json:"GROUP_ID,omitempty"`
	ParentID      int              `json:"PARENT_ID"`
	CustomFields  CustomTaskFields `json:"custom_fields,omitempty"`
}

// Структура для общего тела запроса
type TaskRequest struct {
	Fields interface{} `json:"fields"`
}

type TaskResponse struct {
	Result int `json:"result"`
}

// AddTaskToGroup создает задачу в группе (без привязки к PARENT_ID)
func AddTaskToGroup(title string, responsibleID, groupID int) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := TaskRequest{
		Fields: Task{
			Title:         title,
			ResponsibleID: responsibleID,
			GroupID:       groupID,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response TaskResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if response.Result == 0 {
		return 0, fmt.Errorf("failed to create task, response: %s", string(responseData))
	}

	log.Println("Task created in group with ID:", response.Result)
	return response.Result, nil
}

// AddTaskToParentId создает подзадачу с привязкой к PARENT_ID и пользовательскими полями
func AddTaskToParentId(title string, responsibleID, groupID, parentID int, customFields CustomTaskFields) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Подготовка тела запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"TITLE":          title,
			"RESPONSIBLE_ID": responsibleID,
			"GROUP_ID":       groupID,
			"PARENT_ID":      parentID,
			// Добавляем пользовательские поля
			"UF_AUTO_303168834495": customFields.OrderNumber,       // № заказа
			"UF_AUTO_876283676967": customFields.Customer,          // Заказчик
			"UF_AUTO_794809224848": customFields.Manager,           // Менеджер
			"UF_AUTO_468857876599": customFields.Material,          // Материал
			"UF_AUTO_497907774817": customFields.Comment,           // Комментарий
			"UF_AUTO_433735177517": customFields.ProductionTask,    // Произв. Задача
			"UF_AUTO_726724682983": customFields.Bend,              // Гибка
			"UF_AUTO_512869473370": customFields.Coating,           // Покрытие
			"UF_AUTO_555642596740": customFields.TemporaryOrderSum, // Временная сумма заказа
			"UF_AUTO_552243496167": customFields.Quantity,          // Кол-во
		},
	}

	// Сериализация тела запроса
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	// Разбираем ответ
	var response TaskResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Проверяем успешность создания подзадачи
	if response.Result == 0 {
		return 0, fmt.Errorf("failed to create subtask, response: %s", string(responseData))
	}

	log.Println("Subtask created with ID:", response.Result)
	return response.Result, nil
}

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

func AddMainTasks() {

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
