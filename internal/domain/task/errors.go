package task

import "errors"

var (
	ErrNotFound         = errors.New("task not found")
	ErrInvalidRepeat    = errors.New("invalid repeat rule")
	ErrInvalidDateRange = errors.New("invalid date range")
)
