package payroll

import "errors"

var ErrInvalidDateRange = errors.New("end date cannot be before start date")
