package math_test

import (
	math "broker/internal/math"
	"reflect"
	"testing"
)

func TestReverseSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	math.ReverseSlice(slice)

	expected := []int{5, 4, 3, 2, 1}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestReverseSlicePtr(t *testing.T) {
	a, b, c := 1, 2, 3
	slice := []*int{&a, &b, &c}

	// Expected result after reversal
	expected := []*int{&c, &b, &a}
	math.ReverseSlicePtr(slice)

	got := make([]int, len(slice))
	want := make([]int, len(expected))
	for i, v := range slice {
		got[i] = *v
		want[i] = *expected[i]
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}
