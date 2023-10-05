package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMinuteRangeForTime(t *testing.T) {
	{
		from, to := GetMinuteRangeForTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), from)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), to)
	}
	{
		from, to := GetMinuteRangeForTime(time.Date(2020, 1, 1, 0, 0, 30, 0, time.UTC))
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), from)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), to)
	}
	{
		from, to := GetMinuteRangeForTime(time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), from)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), to)
	}
	{
		from, to := GetMinuteRangeForTime(time.Date(2020, 1, 1, 0, 1, 30, 0, time.UTC))
		assert.Equal(t, time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), from)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), to)
	}
}

func TestAbs(t *testing.T) {
	testCases := []struct {
		name     string
		input    int64
		expected int64
	}{
		{"positive number", 10, 10},
		{"negative number", -10, 10},
		{"zero", 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Abs(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		name     string
		inputA   int64
		inputB   int64
		expected int64
	}{
		{"a is greater than b", 10, 5, 10},
		{"b is greater than a", 5, 10, 10},
		{"a and b are equal", 5, 5, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Max(tc.inputA, tc.inputB)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMin(t *testing.T) {
	testCases := []struct {
		name     string
		inputA   int64
		inputB   int64
		expected int64
	}{
		{"a is less than b", 5, 10, 5},
		{"b is less than a", 10, 5, 5},
		{"a and b are equal", 5, 5, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Min(tc.inputA, tc.inputB)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
