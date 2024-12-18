package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func pullCustomFieldInSmartProcess(checkCoating bool, entityTypeId, smartProcessID int, fieldName, fieldValue string, tasksIDs []int) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.item.update"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	if len(tasksIDs) == 0 {
		return fmt.Errorf("tasksIDs array is empty, cannot update smart process")
	}
	log.Printf("Updating smart process ID %d with tasksIDs: %v", smartProcessID, tasksIDs)

	// Преобразование tasksIDs в []string
	stringTasksIDs := make([]string, len(tasksIDs))
	for i, id := range tasksIDs {
		stringTasksIDs[i] = strconv.Itoa(id)
	}

	if checkCoating == true {

	}
	// Обновляем значение полей в запросе
	requestBody := map[string]interface{}{
		"entityTypeId": entityTypeId,
		"id":           smartProcessID,
		"fields": map[string]interface{}{
			fieldName: stringTasksIDs, // Используем динамическое имя поля
		},
	}
	if checkCoating == true {
		requestBody = map[string]interface{}{
			"entityTypeId": entityTypeId,
			"id":           smartProcessID,
			"fields": map[string]interface{}{
				"ufCrm6_1733264270": "да", // Используем динамическое имя поля
			},
		}
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}
	log.Printf("Request Body: %s", string(jsonData))

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	log.Printf("Response from pullCustomFieldInSmartProcess: %s", string(responseData))

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if _, ok := response["error"]; ok {
		return fmt.Errorf("failed to update smart process: %s", response["error_description"])
	}

	log.Printf("Smart process updated successfully for ID: %d with tasks: %v", smartProcessID, tasksIDs)
	return nil
}
