package services

import "github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"

func validIDs(ids ...int64) error {
	for _, id := range ids {
		if id <= 0 {
			return &resources.ValidationError{Field: "id", Message: "ожидается положительный ID"}
		}
	}
	return nil
}
