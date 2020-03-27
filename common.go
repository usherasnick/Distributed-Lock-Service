package distlock

import "errors"

var (
	ErrDeadlock = errors.New("trying to acquire a lock twice by the same process")
)
