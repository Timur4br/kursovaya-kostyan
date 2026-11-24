package helpers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"api_gateway/entities"
)

func ExtractFilterFromQueryParams(params url.Values) (*entities.Filter, error) {
	config := &entities.Filter{}
	hasValue := false

	if err := setIntParam(params, "camera_id", &config.CameraID); err != nil {
		return nil, err
	}
	if config.CameraID != nil {
		hasValue = true
	}

	if !hasValue {
		return nil, nil
	}
	return config, nil
}

func setStringParam(params url.Values, key string, target **string) error {
	if val := params.Get(key); val != "" {
		*target = &val
	}
	return nil
}
func setIntParam(params url.Values, key string, target **int) error {
	if val := params.Get(key); val != "" {
		num, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid %s parametr: %w", key, err)
		}
		*target = &num
	}
	return nil
}

func setBoolParam(params url.Values, key string, target **bool) error {
	if val := params.Get(key); val != "" {
		num, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid %s parametr: %w", key, err)
		}
		*target = &num
	}
	return nil
}

func setDateParam(params url.Values, key string, target **time.Time) error {
	if val := params.Get(key); val != "" {
		// Список поддерживаемых форматов (в порядке приоритета)
		formats := []string{
			"02.01.2006",
			time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
			"2006-01-02T15:04:05", // Дата и время
			"2006-01-02 15:04:05", // Альтернативный разделитель
			"2006-01-02T15:04",    // Без секунд
			"2006-01-02 15:04",    // Без секунд с пробелом
			"2006-01-02",          // Только дата
			"02.01.2006 15:04:05", // Европейский формат
			time.RFC1123,          // HTTP формат
		}

		var t time.Time
		var err error
		for _, format := range formats {
			t, err = time.Parse(format, val)
			if err == nil {
				*target = &t
				return nil
			}
		}

		return fmt.Errorf("invalid %s parameter: unsupported date format", key)
	}
	return nil
}
