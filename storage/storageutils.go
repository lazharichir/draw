package storage

import "database/sql"

var ErrNoRows = sql.ErrNoRows

func toAnySlice[T any](values []T) []any {
	anyValues := make([]any, len(values))
	for i, value := range values {
		anyValues[i] = value
	}
	return anyValues
}
