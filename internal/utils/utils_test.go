package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	require.Zero(t, Zero[string]())
	require.Zero(t, Zero[int]())
	require.Zero(t, Zero[float32]())
	require.Zero(t, Zero[float64]())
	require.Zero(t, Zero[uint64]())
	require.Zero(t, Zero[[]any]())
}

func TestSame(t *testing.T) {
	// Test with string
	require.True(t, Same("hello", "hello"))
	require.False(t, Same("hello", "world"))

	// Test with int
	require.True(t, Same(42, 42))
	require.False(t, Same(42, 24))

	// Test with float
	require.True(t, Same(3.14, 3.14))
	require.False(t, Same(3.14, 2.71))

	// Test with bool
	require.True(t, Same(true, true))
	require.False(t, Same(true, false))

	// Test with struct
	type Point struct {
		X, Y int
	}
	require.True(t, Same(Point{1, 2}, Point{1, 2}))
	require.False(t, Same(Point{1, 2}, Point{2, 1}))

	// Test with pointer (nil)
	var p1, p2 *int
	require.True(t, Same(p1, p2))

	// Test with pointer (non-nil)
	v1, v2 := 42, 42
	require.False(t, Same(&v1, &v2)) // Different pointers
	p3 := &v1
	require.True(t, Same(&v1, p3)) // Same pointer
}
