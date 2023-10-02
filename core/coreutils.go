package core

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
