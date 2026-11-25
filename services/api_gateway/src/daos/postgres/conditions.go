package pg_dao

import (
	"fmt"
	"strings"

	"api_gateway/entities"
)

func addFilterByCamera(query string, args *[]any, filter *entities.Filter) string {
	if filter == nil {
		return query
	}
	paramCounter := len(*args)
	if filter.CameraID != nil {
		addArgument(&paramCounter, args, &query, *filter.CameraID, " id =")
	}
	if filter.IP != nil {
		addArgument(&paramCounter, args, &query, *filter.IP, " ip =")
	}
	if filter.IsEnabled != nil {
		addArgument(&paramCounter, args, &query, *filter.IsEnabled, " is_enabled =")
	}
	return addFilterByCreatedAt(query, args, filter)
}

func addFilterByObject(query string, args *[]any, filter *entities.Filter) string {
	if filter == nil {
		return query
	}
	paramCounter := len(*args)
	if filter.ObjectID != nil {
		addArgument(&paramCounter, args, &query, *filter.ObjectID, " id =")
	}
	if filter.Type != nil {
		addArgument(&paramCounter, args, &query, *filter.Type, " type =")
	}
	return addFilterByCreatedAt(query, args, filter)
}

func addFilterByEvent(query string, args *[]any, filter *entities.Filter) string {
	if filter == nil {
		return query
	}
	paramCounter := len(*args)
	if filter.CameraID != nil {
		addArgument(&paramCounter, args, &query, *filter.CameraID, " camera_id =")
	}
	if filter.ObjectID != nil {
		addArgument(&paramCounter, args, &query, *filter.ObjectID, " object_id =")
	}
	if filter.Type != nil {
		addArgument(&paramCounter, args, &query, *filter.Type, " type =")
	}
	if filter.ConfidenceMin != nil {
		addArgument(&paramCounter, args, &query, *filter.ConfidenceMin, " confidence >=")
	}
	if filter.ConfidenceMax != nil {
		addArgument(&paramCounter, args, &query, *filter.ConfidenceMax, " confidence <=")
	}
	return addFilterByCreatedAt(query, args, filter)
}

func addFilterByCreatedAt(query string, args *[]any, filter *entities.Filter) string {
	if filter == nil {
		return query
	}
	paramCounter := len(*args)
	if filter.CreatedFrom != nil {
		addArgument(&paramCounter, args, &query, *filter.CreatedFrom, " created_at >=")
	}
	if filter.CreatedTo != nil {
		addArgument(&paramCounter, args, &query, *filter.CreatedTo, " created_at <=")
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
