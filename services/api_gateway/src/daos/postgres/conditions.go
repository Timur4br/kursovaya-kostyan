package pg_dao

import (
	"fmt"
	"strings"

	"api_gateway/entities"
)

func addFilterByZone(query string, args *[]any, filter *entities.Filter) string {
	paramCounter := len(*args)
	if filter != nil {
		if filter.CameraID != nil {
			addArgument(&paramCounter, args, &query, filter.CameraID, " cz.camera_id =")
		}
	}

	return query
}

// Add pagination by next params:
//  1. limit;
//  2. offset.
func addPagination(query string, args *[]any, limit, offset int) string {
	paramCounter := len(*args)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter+1, paramCounter+2)
	*args = append(*args, limit, offset)

	return query
}

func addCondition(queryPtr *string, condition string) {
	if !strings.Contains(*queryPtr, "WHERE") {
		*queryPtr += " WHERE"
	} else {
		*queryPtr += " AND"
	}

	*queryPtr += condition
}

func addArgument(counterPtr *int, argsPtr *[]any, queryPtr *string, param any, argSQL string) {
	*counterPtr++
	addCondition(queryPtr, fmt.Sprintf(argSQL+" $%d", *counterPtr))
	*argsPtr = append(*argsPtr, param)
}
