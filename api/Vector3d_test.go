package api

import (
	"reflect"
	"testing"
)

type testCase struct {
	a, b, expected *Vector3d
}

func TestVector3d_Add(t *testing.T) {
	testCases := []testCase{
		{
			a: &Vector3d{
				X: 1,
				Y: 1,
				Z: 1,
				W: 0,
			},
			b: &Vector3d{
				X: 1,
				Y: 1,
				Z: 1,
				W: 1,
			},
			expected: &Vector3d{
				X: 2,
				Y: 2,
				Z: 2,
				W: 0,
			},
		},
	}

	for _, test := range testCases {
		result := test.a.Add(test.b)
		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("Add: expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector3d_Sub(t *testing.T) {
	testCases := []testCase{
		{
			a: &Vector3d{
				X: 2,
				Y: 2,
				Z: 2,
				W: 0,
			},
			b: &Vector3d{
				X: 1,
				Y: 1,
				Z: 1,
				W: 1,
			},
			expected: &Vector3d{
				X: 1,
				Y: 1,
				Z: 1,
				W: 0,
			},
		},
	}

	for _, test := range testCases {
		result := test.a.Sub(test.b)
		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("Sub: expected %v, got %v", test.expected, result)
		}
	}
}
