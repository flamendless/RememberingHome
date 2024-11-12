package common

import "errors"

var (
	ERR_QUIT         = errors.New("Exit game has been triggered")
	ERR_NOT_YET_IMPL = errors.New("Not yet implemented")
)
