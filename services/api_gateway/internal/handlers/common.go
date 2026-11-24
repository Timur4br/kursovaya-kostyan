package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const MaxBodyBytes int64 = 1 << 20

type requestError struct {
	status        int
	code, message string
}

func (e *requestError) Error() string { return e.code }

func bindJSON[T any](c *gin.Context) (T, bool) {
	var input T
	media, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || media != "application/json" {
		reject(c, 415, "unsupported_media_type", "Нужен Content-Type: application/json", nil)
		return input, false
	}
	raw, err := readObject(c)
	if err != nil {
		var p *requestError
		if errors.As(err, &p) {
			reject(c, p.status, p.code, p.message, nil)
		} else {
			reject(c, 400, "invalid_json", "Некорректный JSON", nil)
		}
		return input, false
	}
	data, err := json.Marshal(raw)
	if err != nil {
		reject(c, 400, "invalid_json", "Некорректный JSON", nil)
		return input, false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		reject(c, 422, "validation_error", "Проверьте имена и типы полей запроса", nil)
		return input, false
	}
	return input, true
}

func readObject(c *gin.Context) (map[string]json.RawMessage, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	defer c.Request.Body.Close()
	decoder := json.NewDecoder(c.Request.Body)
	token, err := decoder.Token()
	if err != nil {
		return nil, jsonError(err)
	}
	if token != json.Delim('{') {
		return nil, &requestError{400, "invalid_json", "Ожидается JSON-объект"}
	}
	raw := map[string]json.RawMessage{}
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, jsonError(err)
		}
		name, ok := key.(string)
		if !ok {
			return nil, &requestError{400, "invalid_json", "Некорректный JSON"}
		}
		if _, exists := raw[name]; exists {
			return nil, &requestError{400, "duplicate_field", "Поле повторяется в JSON"}
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, jsonError(err)
		}
		raw[name] = value
	}
	if _, err := decoder.Token(); err != nil {
		return nil, jsonError(err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, jsonError(err)
		}
		return nil, &requestError{400, "invalid_json", "После объекта есть лишний JSON"}
	}
	return raw, nil
}

func jsonError(err error) error {
	var large *http.MaxBytesError
	if errors.As(err, &large) {
		return &requestError{413, "body_too_large", "Тело запроса превышает 1 МиБ"}
	}
	return &requestError{400, "invalid_json", "Некорректный JSON"}
}

func pathID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		reject(c, 400, "invalid_id", "ID должен быть положительным целым числом", nil)
		return 0, false
	}
	return id, true
}

func reject(c *gin.Context, status int, code, message string, fields map[string]string) {
	value := gin.H{"code": code, "message": message}
	if len(fields) > 0 {
		value["fields"] = fields
	}
	c.AbortWithStatusJSON(status, gin.H{"error": value})
}

func writeError(c *gin.Context, err error, logger *slog.Logger) {
	status, code, message := 500, "internal_error", "Ошибка сервера"
	var fields map[string]string
	var validation *resources.ValidationError
	var pg *pgconn.PgError
	switch {
	case errors.As(err, &validation):
		status, code, message = 422, "validation_error", "Проверьте поля запроса"
		fields = map[string]string{validation.Field: validation.Message}
	case errors.Is(err, sql.ErrNoRows):
		status, code, message = 404, "not_found", "Запись не найдена"
	case errors.Is(err, context.DeadlineExceeded):
		status, code, message = 504, "timeout", "Время обработки запроса истекло"
	case errors.Is(err, context.Canceled):
		status, code, message = 503, "request_canceled", "Запрос отменён"
	case errors.As(err, &pg):
		switch pg.Code {
		case "23505":
			status, code, message = 409, "already_exists", "Запись с такими уникальными значениями уже существует"
		case "23503":
			status, code, message = 409, "relation_conflict", "Связанная запись отсутствует или используется другими данными"
		case "23514", "23502":
			status, code, message = 422, "constraint_violation", "Значения нарушают ограничения данных"
		case "57014":
			status, code, message = 504, "timeout", "Время обработки запроса истекло"
		default:
			if strings.HasPrefix(pg.Code, "22") {
				status, code, message = 422, "constraint_violation", "Значения нарушают ограничения данных"
			}
			if strings.HasPrefix(pg.Code, "08") {
				status, code, message = 503, "database_unavailable", "База данных недоступна"
			}
		}
	}
	if status >= 500 {
		logger.Error("entity operation failed", "code", code, "error_type", fmt.Sprintf("%T", err))
	}
	reject(c, status, code, message, fields)
}
