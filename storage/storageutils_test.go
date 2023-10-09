package storage

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNoRows(t *testing.T) {
	assert.True(t, errors.Is(ErrNoRows, sql.ErrNoRows))
	assert.True(t, ErrNoRows == sql.ErrNoRows)
}
