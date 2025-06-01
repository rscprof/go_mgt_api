// Package api содержит тесты для клиента API транспортных данных Москвы.
// Этот файл включает unit-тесты для проверки корректности работы метода GetStopData.
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testStopID - тестовый идентификатор остановки для проверки.
// Соответствует UUID-формату, используемому в API.
const testStopID = "9d7f733a-d532-4fca-a922-4c978b79681c"

// fakeAPIResponse - фиктивный JSON-ответ API, используемый для тестирования.
// Включает все возможные поля ответа с тестовыми значениями.
const fakeAPIResponse = `
{
    "id": "9d7f733a-d532-4fca-a922-4c978b79681c",
    "name": "ул. Льва Толстого",
    "type": "ground",
    "routePath": [
        {
            "id": "bus_123",
            "type": "bus",
            "number": "М10",
            "lastStopName": "Киевский вокзал",
            "color": "#FF0000",
            "fontColor": "#FFFFFF",
            "externalForecast": [
                {
                    "time": 1719752400,
                    "byTelemetry": 1,
                    "tmId": 987654,
                    "routePathId": "bus_123"
                }
            ]
        }
    ]
}
`

// TestGetStopData_Success проверяет успешный сценарий получения данных об остановке.
// Тест покрывает:
// - Создание мок-сервера API
// - Подмену базового URL клиента
// - Корректность обработки ответа API
// - Полноту и правильность данных в возвращаемой структуре
func TestGetStopData_Success(t *testing.T) {
	// Arrange (подготовка)
	// Создаем мок-сервер, который возвращает фиксированный JSON-ответ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем что запрос пришел на правильный endpoint
		assert.Contains(t, r.URL.Path, "/stop_v2/"+testStopID, "Неверный endpoint API")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fakeAPIResponse))
	}))
	defer server.Close()

	// Создаем клиент API и подменяем базовый URL на адрес мок-сервера
	client := NewClient()
	client.SetBaseURL(server.URL + "/")

	// Act (действие)
	stopData, err := client.GetStopData(testStopID)

	// Assert (проверки)
	// Проверяем что не было ошибок
	assert.NoError(t, err, "Неожиданная ошибка при вызове GetStopData")
	assert.NotNil(t, stopData, "Данные об остановке не должны быть nil")

	// Проверяем основные поля остановки
	assert.Equal(t, testStopID, stopData.ID, "ID остановки не совпадает")
	assert.Equal(t, "ул. Льва Толстого", stopData.Name, "Название остановки не совпадает")
	assert.Equal(t, "ground", stopData.Type, "Тип остановки не совпадает")

	// Проверяем что есть ровно один маршрут
	assert.Len(t, stopData.RoutePath, 1, "Должен быть ровно один маршрут")
	route := stopData.RoutePath[0]

	// Проверяем данные маршрута
	assert.Equal(t, "bus_123", route.ID, "ID маршрута не совпадает")
	assert.Equal(t, "bus", route.Type, "Тип транспорта не совпадает")
	assert.Equal(t, "М10", route.Number, "Номер маршрута не совпадает")
	assert.Equal(t, "Киевский вокзал", route.LastStopName, "Название конечной остановки не совпадает")
	assert.Equal(t, "#FF0000", route.Color, "Цвет маршрута не совпадает")
	assert.Equal(t, "#FFFFFF", route.FontColor, "Цвет текста маршрута не совпадает")

	// Проверяем прогноз прибытия
	assert.Len(t, route.ExternalForecast, 1, "Должен быть ровно один прогноз прибытия")
	forecast := route.ExternalForecast[0]

	assert.Equal(t, 1719752400, forecast.Time, "Время прибытия не совпадает")
	assert.Equal(t, 1, forecast.ByTelemetry, "Флаг телеметрии не совпадает")
	assert.Equal(t, 987654, forecast.TmID, "ID транспорта не совпадает")
	assert.Equal(t, "bus_123", forecast.RoutePathID, "ID маршрута в прогнозе не совпадает")
}
