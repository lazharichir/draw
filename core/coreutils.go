package core

import "time"

func GetMinuteRangeForTime(t time.Time) (time.Time, time.Time) {
	from := t.Truncate(time.Minute)
	to := from.Add(time.Minute)
	return from, to
}

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
