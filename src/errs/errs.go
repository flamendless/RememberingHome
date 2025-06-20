package errs

import "errors"

var (
	ERR_QUIT         = errors.New("Exit game has been triggered")
	ERR_NIL_ELEM     = errors.New("Nil element in array")
	ERR_NOT_YET_IMPL = errors.New("Not yet implemented")
)
