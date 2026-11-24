package helpers

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"api_gateway/entities"
)

func ExtractFilterFromQueryParams(params url.Values) (*entities.Filter, error) {
	config := &entities.Filter{}

	if err := setIntParam(params, "camera_id", &config.CameraID); err != nil {
		return nil, err
	}
	if err := setIntParam(params, "object_id", &config.ObjectID); err != nil {
		return nil, err
	}
	if err := setStringParam(params, "type", &config.Type); err != nil {
		return nil, err
	}
	if err := setStringParam(params, "ip", &config.IP); err != nil {
		return nil, err
	}
	if err := setBoolParam(params, "is_enabled", &config.IsEnabled); err != nil {
		return nil, err
	}
	if err := setFloatParam(params, "confidence_min", &config.ConfidenceMin); err != nil {
		return nil, err
	}
	if err := setFloatParam(params, "confidence_max", &config.ConfidenceMax); err != nil {
		return nil, err
	}
	if err := setDateParam(params, "created_from", &config.CreatedFrom); err != nil {
		return nil, err
	}
	if err := setDateParam(params, "created_to", &config.CreatedTo); err != nil {
		return nil, err
	}

	if config.ConfidenceMin != nil && config.ConfidenceMax != nil && *config.ConfidenceMin > *config.ConfidenceMax {
		return nil, fmt.Errorf("confidence_min must not exceed confidence_max")
	}
	if config.CreatedFrom != nil && config.CreatedTo != nil && config.CreatedFrom.After(*config.CreatedTo) {
		return nil, fmt.Errorf("created_from must not exceed created_to")
	}
	if config.CameraID == nil && config.ObjectID == nil && config.Type == nil && config.IP == nil &&
		config.IsEnabled == nil && config.ConfidenceMin == nil && config.ConfidenceMax == nil &&
		config.CreatedFrom == nil && config.CreatedTo == nil {
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

func setFloatParam(params url.Values, key string, target **float64) error {
	if val := params.Get(key); val != "" {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("invalid %s parameter: %w", key, err)
		}
		if math.IsNaN(num) || math.IsInf(num, 0) {
			return fmt.Errorf("invalid %s parameter: value must be finite", key)
		}
		*target = &num
	}
	return nil
}
