// Package api предоставляет клиент для работы с API транспортных данных Москвы.
// Включает методы для получения информации об остановках и маршрутах общественного транспорта.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Константы для настройки клиента
const (
	// DEFAULT_BASE_API_URL - базовый URL API транспортных данных Москвы
	DEFAULT_BASE_API_URL = "https://api.moscowapp.mos.ru/v8.2/"

	// USER_AGENT_HEADER - User-Agent для HTTP-запросов
	USER_AGENT_HEADER = "Mozilla/5.0"

	// ACCEPT_HEADER - Accept-заголовок для HTTP-запросов
	ACCEPT_HEADER = "application/json"

	// DEBUG - флаг отладки (вывод дополнительной информации)
	DEBUG = true
)

// APIClient представляет клиент для работы с Moscow API.
// Содержит настройки подключения и методы для выполнения запросов.
type APIClient struct {
	baseURL string       // Базовый URL API
	client  *http.Client // HTTP-клиент с таймаутом
}

// NewClient создает новый экземпляр APIClient с настройками по умолчанию.
// Возвращает:
//
//	*APIClient - инициализированный клиент API
//
// Пример:
//
//	client := api.NewClient()
//	data, err := client.GetStopData("stop_id")
func NewClient() *APIClient {
	return &APIClient{
		baseURL: DEFAULT_BASE_API_URL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// SetBaseURL устанавливает кастомный базовый URL API.
// Используется преимущественно для тестирования.
//
// Параметры:
//
//	url - новый базовый URL (например, "http://localhost:8080/api/")
func (c *APIClient) SetBaseURL(url string) {
	c.baseURL = url
}

// doRequest выполняет GET-запрос к указанному эндпоинту API.
// Внутренний метод, используется другими методами клиента.
//
// Параметры:
//
//	endpoint - относительный путь эндпоинта (например, "stop_v2/id")
//
// Возвращает:
//
//	[]byte - тело ответа
//	error - ошибка, если запрос не удался
//
// Возможные ошибки:
//   - ошибка создания запроса
//   - ошибка сети
//   - статус код != 200
func (c *APIClient) doRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", USER_AGENT_HEADER)
	req.Header.Set("Accept", ACCEPT_HEADER)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if DEBUG {
			fmt.Printf("DEBUG: Error response body: %s\n", string(body))
		}
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GetStopData получает данные об остановке по её идентификатору.
// Реализует метод APIClientInterface.
//
// Параметры:
//
//	stopID - уникальный идентификатор остановки (формат UUID)
//
// Возвращает:
//
//	*StopData - структура с данными об остановке
//	error - ошибка, если запрос или парсинг не удались
//
// Пример:
//
//	data, err := client.GetStopData("9d7f733a-d532-4fca-a922-4c978b79681c")
//	if err != nil {
//	    // обработка ошибки
//	}
//	fmt.Println("Название остановки:", data.Name)
func (c *APIClient) GetStopData(stopID string) (*StopData, error) {
	body, err := c.doRequest("stop_v2/" + stopID)
	if err != nil {
		return nil, err
	}

	var data StopData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	return &data, nil
}
