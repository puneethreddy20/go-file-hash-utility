package main

import (
	"testing"
)

func TestIsThreadsnumValid(t *testing.T) {
	thread := []struct {
		x      int
		output bool
	}{
		{0, false}, {10, true}, {8, true}, {24, true}, {-1, false},
	}
	for _, value := range thread {
		result, _ := IsThreadsnumValid(value.x)
		if result != value.output {
			t.Errorf("expected %t but got %t", value.output, result)
		}
	}

}

func TestIsWindowValid(t *testing.T) {
	window := []struct {
		x      int
		output bool
	}{
		{0, false}, {10, false}, {8, true}, {30, false}, {24, true}, {-1, false},
	}
	for _, value := range window {
		result, _ := IsWindowValid(value.x)
		if result != value.output {
			t.Errorf("expected %t but got %t", value.output, result)
		}
	}

}

func TestAddPadding(t *testing.T) {
	inputcases := []struct {
		s      string
		size   int
		strlen int
	}{
		{"abcdef", 8, 8}, {"abc1234abc123", 16, 16}, {"abc1234abc123", 8, 0},
	}
	for _, value := range inputcases {
		result, _ := AddPadding(value.s, value.size)
		if len(result) != value.strlen {
			t.Errorf("expected length %d but got length %d", value.strlen, len(result))
		}
	}
}

func TestConvertchunckstrtouint64(t *testing.T) {
	inputcases := []struct {
		s      string
		output uint64
	}{
		{"A", 65}, {"ab", 24930}, {"abc1234abc123", 3689681069922136627}, {"ABCD", 1094861636},
	}
	for _, value := range inputcases {
		result, _ := Convertchunckstrtouint64(value.s)
		if result != value.output {
			t.Errorf("expected length %d but got length %d", value.output, result)
		}
	}
}

func TestGetAllchuncks(t *testing.T) {
	inputcases := []struct {
		s      string
		size   int
		length int
	}{
		{"abcdefghij", 8, 2}, {"abc1234abc123", 16, 1}, {"abc1234abc123", 8, 2},
	}
	for _, value := range inputcases {
		result, _ := GetAllchuncks(value.s, value.size)
		if len(result) != value.length {
			t.Errorf("expected length %d but got length %d", value.length, len(result))
		}
	}
}
