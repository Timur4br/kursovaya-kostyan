package httpapi

import "github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"

type schemaMap = map[string]any

func reference(name string) schemaMap { return schemaMap{"$ref": "#/components/schemas/" + name} }
func content(schema any) schemaMap    { return schemaMap{"application/json": schemaMap{"schema": schema}} }
func fieldSchema(f resources.Field) schemaMap {
	s := schemaMap{}
	switch f.Kind {
	case resources.ID:
		s["type"] = "integer"
		s["format"] = "int64"
		s["minimum"] = 1
	case resources.Number:
		s["type"] = "number"
		s["format"] = "double"
		if f.Positive {
			s["minimum"] = 0
			s["exclusiveMinimum"] = true
		}
		if f.Confidence {
			s["minimum"] = 0
			s["maximum"] = 1
		}
	case resources.Timestamp:
		s["type"] = "string"
		s["format"] = "date-time"
	case resources.Object:
		s["type"] = "object"
		s["additionalProperties"] = true
	default:
		s["type"] = "string"
		if f.MaxLength > 0 {
			s["maxLength"] = f.MaxLength
		}
		if f.Required {
			s["minLength"] = 1
		}
	}
	if len(f.Enum) > 0 {
		s["enum"] = f.Enum
	}
	if f.Nullable {
		s["nullable"] = true
	}
	if f.ReadOnly {
		s["readOnly"] = true
	}
	return s
}
func OpenAPI() schemaMap {
	schemas := schemaMap{"Error": schemaMap{"type": "object", "required": []string{"error"}, "properties": schemaMap{"error": schemaMap{"type": "object", "required": []string{"code", "message"}, "properties": schemaMap{"code": schemaMap{"type": "string"}, "message": schemaMap{"type": "string"}, "fields": schemaMap{"type": "object", "additionalProperties": schemaMap{"type": "string"}}}}}}}
	paths := schemaMap{}
	for _, r := range resources.All {
		readProps, inputProps := schemaMap{}, schemaMap{}
		required := []string{}
		for _, f := range r.Fields {
			readProps[f.Name] = fieldSchema(f)
			if !f.ReadOnly {
				inputProps[f.Name] = fieldSchema(f)
				if f.Required {
					required = append(required, f.Name)
				}
			}
		}
		schemas[r.Path] = schemaMap{"type": "object", "additionalProperties": false, "properties": readProps}
		input := schemaMap{"type": "object", "additionalProperties": false, "properties": inputProps}
		if len(required) > 0 {
			input["required"] = required
		}
		schemas[r.Path+"_input"] = input
		body := schemaMap{"required": true, "content": content(reference(r.Path + "_input"))}
		record := schemaMap{"type": "object", "required": []string{"data"}, "properties": schemaMap{"data": reference(r.Path)}}
		responses := func(code string) schemaMap {
			result := schemaMap{code: schemaMap{"description": "Успешный результат", "content": content(record)}}
			if code == "204" {
				result[code] = schemaMap{"description": "Удалено; тело ответа отсутствует"}
			}
			for _, status := range []string{"400", "404", "409", "413", "415", "422", "500", "503", "504"} {
				result[status] = schemaMap{"description": "Ошибка запроса или сервера", "content": content(reference("Error"))}
			}
			return result
		}
		paths["/"+r.Path] = schemaMap{"post": schemaMap{"tags": []string{r.Path}, "operationId": "create_" + r.Path, "requestBody": body, "responses": responses("201")}}
		parameters := []any{}
		for _, key := range r.Keys {
			parameters = append(parameters, schemaMap{"name": key, "in": "path", "required": true, "schema": schemaMap{"type": "integer", "format": "int64", "minimum": 1}})
		}
		paths[r.ItemPath()] = schemaMap{
			"parameters": parameters,
			"get":        schemaMap{"tags": []string{r.Path}, "operationId": "get_" + r.Path, "responses": responses("200")},
			"put":        schemaMap{"tags": []string{r.Path}, "operationId": "update_" + r.Path, "description": "Полная замена изменяемых полей; пропущенные необязательные поля сбрасываются к значениям по умолчанию.", "requestBody": body, "responses": responses("200")},
			"delete":     schemaMap{"tags": []string{r.Path}, "operationId": "delete_" + r.Path, "responses": responses("204")},
		}
	}
	return schemaMap{"openapi": "3.0.3", "info": schemaMap{"title": "Курсовая Костян", "version": "0.2.0", "description": "Локальный API управления данными. Авторизация ещё не реализована."}, "paths": paths, "components": schemaMap{"schemas": schemas}}
}
