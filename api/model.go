// Package api предоставляет модели данных для работы с API транспортных данных Москвы.
// Модели соответствуют структуре ответов API и используются для JSON-сериализации.
package api

// Forecast представляет информацию о прогнозе прибытия транспорта на остановку.
type Forecast struct {
	// Time - временная метка прибытия (unix timestamp)
	Time int `json:"time"`

	// ByTelemetry - флаг, указывающий получен ли прогноз по телеметрии (1 - да, 0 - нет)
	ByTelemetry int `json:"byTelemetry"`

	// TmID - идентификатор транспортного средства в системе телеметрии
	TmID int `json:"tmId"`

	// RoutePathID - идентификатор маршрута
	RoutePathID string `json:"routePathId"`
}

// RoutePath описывает маршрут общественного транспорта, проходящий через остановку.
type RoutePath struct {
	// ID - уникальный идентификатор маршрута
	ID string `json:"id"`

	// Type - тип транспорта (например, "bus", "trolley", "tram")
	Type string `json:"type"`

	// Number - номер маршрута (например, "27", "М10")
	Number string `json:"number"`

	// LastStopName - название конечной остановки маршрута
	LastStopName string `json:"lastStopName"`

	// Color - цвет маршрута в HEX-формате (например, "#FF0000")
	Color string `json:"color"`

	// FontColor - цвет текста для отображения номера маршрута
	FontColor string `json:"fontColor"`

	// ExternalForecast - список прогнозов прибытия транспорта
	ExternalForecast []Forecast `json:"externalForecast"`
}

// StopData содержит полную информацию об остановке общественного транспорта.
// Возвращается методом GetStopData APIClient.
type StopData struct {
	// ID - уникальный идентификатор остановки (UUID)
	ID string `json:"id"`

	// Name - название остановки (например, "ул. Льва Толстого")
	Name string `json:"name"`

	// Type - тип остановки (например, "ground" - наземный транспорт)
	Type string `json:"type"`

	// RoutePath - список маршрутов, проходящих через остановку
	RoutePath []RoutePath `json:"routePath"`
}
