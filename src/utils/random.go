package utils

import "math/rand/v2"

func IntRandRange(x int, y int) int {
	return rand.IntN(y-x) + x
}
