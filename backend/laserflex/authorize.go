package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Request struct {
	AuthID           string
	AuthExpires      int
	RefreshID        string
	MemberID         string
	Status           string
	Placement        string
	PlacementOptions string
}

var (
	GlobalAuthID    string
	GlobalRefreshID string
	GlobalMemberID  string
	AuthExpiryTime  time.Time
	UpdateInterval  = 30 * time.Second // Обновление каждые 50 минут
	ClientID        = "local.671aad770da2d1.64572237"
	ClientSecret    = "qWQFm8UJThmJQNl6BfzVjhgfhlCFALKSG586NBD1zjDcdn8ISG"
	OAuthURL        = "https://oauth.bitrix.info/oauth/token/"
)

func AuthorizeEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("resp_at_first:", string(bs))

	authValues := ParseValuesLaserflex(w, bs)
	if authValues.AuthID == "" || authValues.RefreshID == "" {
		http.Error(w, "Missing required auth data", http.StatusBadRequest)
		return
	}

	// Сохраняем глобальные значения
	GlobalAuthID = authValues.AuthID
	GlobalRefreshID = authValues.RefreshID
	GlobalMemberID = authValues.MemberID
	AuthExpiryTime = time.Now().Add(time.Duration(authValues.AuthExpires) * time.Second)

	// Запускаем автоматическое обновление токена
	go StartTokenRefresh()

	log.Printf("Authorization successful: AuthID=%s, RefreshID=%s, MemberID=%s", GlobalAuthID, GlobalRefreshID, GlobalMemberID)
}

func ParseValuesLaserflex(w http.ResponseWriter, bs []byte) Request {
	values, err := url.ParseQuery(string(bs))
	if err != nil {
		log.Println("error parsing query:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return Request{}
	}

	authExpires, err := strconv.Atoi(values.Get("AUTH_EXPIRES"))
	if err != nil {
		log.Println("error converting AUTH_EXPIRES to int:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return Request{}
	}

	return Request{
		AuthID:           values.Get("AUTH_ID"),
		AuthExpires:      authExpires,
		RefreshID:        values.Get("REFRESH_ID"),
		MemberID:         values.Get("member_id"),
		Status:           values.Get("status"),
		Placement:        values.Get("PLACEMENT"),
		PlacementOptions: values.Get("PLACEMENT_OPTIONS"),
	}
}

func RefreshToken() error {
	// Проверяем, есть ли refresh_token
	if GlobalRefreshID == "" {
		log.Println("RefreshToken is empty")
		return fmt.Errorf("refresh_token is empty")
	}

	// Формируем URL для обновления токена
	serverEndpoint := "https://oauth.bitrix.info/oauth/token/" // Убедитесь, что это корректный URL
	if serverEndpoint == "" {
		log.Println("Server endpoint is empty")
		return fmt.Errorf("server endpoint is empty")
	}

	// Формируем тело запроса
	requestBody, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     ClientID,
		"client_secret": ClientSecret,
		"refresh_token": GlobalRefreshID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal refresh token request: %w", err)
	}

	// Выполняем запрос на обновление токена
	resp, err := http.Post(serverEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to send refresh token request: %w", err)
	}
	defer resp.Body.Close()

	// Логирование ответа
	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("RefreshToken Response: %s", string(responseBody))

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token, status: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	// Читаем и парсим ответ
	var response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return fmt.Errorf("failed to decode refresh token response: %w", err)
	}

	// Проверяем, получили ли мы токены
	if response.AccessToken == "" || response.RefreshToken == "" {
		log.Println("Received empty tokens from Bitrix24")
		return fmt.Errorf("received empty tokens from Bitrix24")
	}

	// Обновляем глобальные значения
	GlobalAuthID = response.AccessToken
	GlobalRefreshID = response.RefreshToken
	AuthExpiryTime = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

	log.Printf("Token refreshed successfully: AuthID=%s, RefreshID=%s", GlobalAuthID, GlobalRefreshID)
	return nil
}

func StartTokenRefresh() {
	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := RefreshToken(); err != nil {
			log.Println("Error refreshing token:", err)
		}
	}
}
