package utils

func ClampInt(value, minValue, maxValue int) int {
	return max(minValue, min(value, maxValue))
}

func ClampFloat64(value, minValue, maxValue float64) float64 {
	return max(minValue, min(value, maxValue))
}

func WrapAroundInt(value, low, high int) int {
	return (value-low)%(high-low) + low
}
