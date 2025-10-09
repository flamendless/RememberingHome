package errs

import "errors"

var (
	ErrQuit       = errors.New("exit game has been triggered")
	ErrNilElem    = errors.New("nil element in array")
	ErrNotYetImpl = errors.New("not yet implemented")
)
