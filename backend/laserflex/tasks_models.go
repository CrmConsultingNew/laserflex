package laserflex

import "encoding/json"

// Структуры для задачи
type Task struct {
	Title         string   `json:"TITLE"`
	ResponsibleID int      `json:"RESPONSIBLE_ID"`
	GroupID       int      `json:"GROUP_ID,omitempty"`
	UfCrmTask     []string `json:"UF_CRM_TASK,omitempty"`
}

// Структура для пользовательских полей задачи
type CustomTaskFields struct {
	Quantity          string `json:"UF_AUTO_552243496167,omitempty"` // Кол-во
	TemporaryOrderSum string `json:"UF_AUTO_555642596740,omitempty"` // Временная сумма заказа
	OrderNumber       string `json:"UF_AUTO_303168834495,omitempty"` // № заказа
	Customer          string `json:"UF_AUTO_876283676967,omitempty"` // Заказчик
	Manager           string `json:"UF_AUTO_794809224848,omitempty"` // Менеджер
	Material          string `json:"UF_AUTO_468857876599,omitempty"` // Материал
	Bend              string `json:"UF_AUTO_726724682983,omitempty"` // Гибка
	ProductionTask    string `json:"UF_AUTO_433735177517,omitempty"` // Произв. Задача
	Comment           string `json:"UF_AUTO_497907774817,omitempty"` // Комментарий
	Coating           string `json:"UF_AUTO_512869473370,omitempty"` // Покрытие
	AllowTimeTracking string `json:"ALLOW_TIME_TRACKING,omitempty"`
	TimeEstimate      int    `json:"TIME_ESTIMATE,omitempty"` //
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
	Fields map[string]interface{} `json:"fields"`
}

type TaskResponse struct {
	Result struct {
		Task struct {
			ID json.RawMessage `json:"id"`
		} `json:"task"`
	} `json:"result"`
}
