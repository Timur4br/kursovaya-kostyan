package helpers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultLimit  int64 = 10
	defaultOffset int64 = 0
)

func ValidatePagination(limit int64, offset int64) error {
	if limit < 0 || offset < 0 {
		return fmt.Errorf("wrong pagination parameters value")
	}

	return nil
}

func ExtractPagination(r *http.Request) (int64, int64, error) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	limitInt64, err := strconv.ParseInt(params.Get("limit"), 10, 64)
	if err != nil {
		return defaultLimit, defaultOffset, fmt.Errorf("wrong limit value format")
	}

	offsetInt64, err := strconv.ParseInt(params.Get("offset"), 10, 64)
	if err != nil {
		return defaultLimit, defaultOffset, fmt.Errorf("wrong offset value format")
	}

	return limitInt64, offsetInt64, nil
}
