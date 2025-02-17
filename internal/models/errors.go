package models

import (
	"errors"
)

var ErrNoRecord = errors.New("sql: no rows in result set")
