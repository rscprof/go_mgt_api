// Package api содержит интеграционные тесты для клиента API транспортных данных Москвы.
// Тесты в этом файле выполняют реальные запросы к внешнему API и требуют сетевого подключения.
package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// realStopID содержит идентификатор реальной остановки для тестирования.
// Это остановка "ул. Льва Толстого" в Москве.
const realStopID = "9d7f733a-d532-4fca-a922-4c978b79681c"

// setEnvForIntegrationTest временно активирует интеграционные тесты через переменную окружения.
// Возвращает функцию для восстановления исходного состояния окружения.
//
// Пример использования:
//
//	restore := setEnvForIntegrationTest()
//	defer restore()
func setEnvForIntegrationTest() (restore func()) {
	initialValue := os.Getenv("RUN_INTEGRATION_TESTS")
	os.Setenv("RUN_INTEGRATION_TESTS", "true")
	return func() { os.Setenv("RUN_INTEGRATION_TESTS", initialValue) }
}

// isIntegrationTestEnabled проверяет, разрешено ли выполнение интеграционных тестов.
// Тесты выполняются только при RUN_INTEGRATION_TESTS=true.
func isIntegrationTestEnabled() bool {
	return os.Getenv("RUN_INTEGRATION_TESTS") == "true"
}

// TestGetStopData_RealAPI проверяет работу клиента с реальным API транспортных данных.
// Это интеграционный тест, который:
// 1. Выполняет реальный HTTP-запрос к API Москвы
// 2. Проверяет корректность полученных данных
// 3. Требует сетевого подключения
//
// Для запуска:
//
//	RUN_INTEGRATION_TESTS=true go test -v ./... -run TestGetStopData_RealAPI
func TestGetStopData_RealAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	if !isIntegrationTestEnabled() {
		t.Skip("Интеграционные тесты отключены. Для запуска установите RUN_INTEGRATION_TESTS=true")
	}

	// Создаем клиент с дефолтными настройками (реальный API URL)
	client := NewClient()

	// Выполняем запрос к API
	stopData, err := client.GetStopData(realStopID)

	// Базовые проверки
	require.NoError(t, err, "Ошибка при запросе данных об остановке")
	require.NotNil(t, stopData, "Данные об остановке не должны быть nil")

	t.Run("Проверка основных полей остановки", func(t *testing.T) {
		assert.Equal(t, realStopID, stopData.ID, "Неверный ID остановки")
		assert.NotEmpty(t, stopData.Name, "Название остановки должно быть заполнено")
		assert.NotEmpty(t, stopData.Type, "Тип остановки должен быть указан")
	})

	if len(stopData.RoutePath) > 0 {
		t.Run("Проверка данных маршрута", func(t *testing.T) {
			route := stopData.RoutePath[0]
			assert.NotEmpty(t, route.Number, "Номер маршрута должен быть указан")
			assert.NotEmpty(t, route.LastStopName, "Название конечной остановки должно быть указано")

			if len(route.ExternalForecast) > 0 {
				forecast := route.ExternalForecast[0]
				assert.Positive(t, forecast.Time, "Время прибытия должно быть положительным")
				assert.NotEmpty(t, forecast.RoutePathID, "ID маршрута в прогнозе должен быть указан")
			}
		})
	}

	// Дополнительные проверки можно добавить здесь
	t.Logf("Успешно получены данные по остановке: %s (%s)", stopData.Name, stopData.ID)
}
