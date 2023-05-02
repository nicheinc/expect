package expect

import (
	"testing"
)

func TestEqual(t *testing.T) {
	var (
		one      = 1
		onePrime = 1
		two      = 2
	)
	testCases := []struct {
		name   string
		first  *int
		second *int
		fail   bool
	}{
		{
			name:   "NotEqual",
			first:  &one,
			second: &two,
			fail:   true,
		},
		{
			name:   "SameReferents",
			first:  &one,
			second: &onePrime,
			fail:   false,
		},
		{
			name:   "SamePointer",
			first:  &one,
			second: &one,
			fail:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tt := new(testing.T)
			Equal(tt, testCase.first, testCase.second)
			fail := tt.Failed()
			if testCase.fail != fail {
				t.Errorf("Expected failure: %v\nActual failure: %v", testCase.fail, fail)
			}
		})
	}
}

func TestEqualUnordered(t *testing.T) {
	testCases := []struct {
		name   string
		first  sliceWrapper
		second sliceWrapper
		fail   bool
	}{
		{
			name:   "Nil/Nil",
			first:  nil,
			second: nil,
			fail:   false,
		},
		{
			name:   "Nil/Empty",
			first:  nil,
			second: []intWrapper{},
			fail:   true,
		},
		{
			name:   "Nil/NonNil",
			first:  nil,
			second: []intWrapper{1},
			fail:   true,
		},
		{
			name:   "DifferentElements",
			first:  []intWrapper{1, 2, 3},
			second: []intWrapper{4, 5, 6},
			fail:   true,
		},
		{
			name:   "SameElements/SameOrder",
			first:  []intWrapper{1, 2, 3},
			second: []intWrapper{1, 2, 3},
			fail:   false,
		},
		{
			name:   "SameElements/DifferentOrder",
			first:  []intWrapper{1, 2, 3},
			second: []intWrapper{3, 1, 2},
			fail:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tt := new(testing.T)
			EqualUnordered(tt, testCase.first, testCase.second)
			fail := tt.Failed()
			if testCase.fail != fail {
				t.Errorf("Expected failure: %v\nActual failure: %v", testCase.fail, fail)
			}
		})
	}
}

type intWrapper int

type sliceWrapper []intWrapper

func TestDeepEqual(t *testing.T) {
	var (
		one      = 1
		onePrime = 1
		two      = 2
	)
	testCases := []struct {
		name   string
		first  *int
		second *int
		fail   bool
	}{
		{
			name:   "NotEqual",
			first:  &one,
			second: &two,
			fail:   true,
		},
		{
			name:   "SameReferents",
			first:  &one,
			second: &onePrime,
			fail:   false,
		},
		{
			name:   "SamePointer",
			first:  &one,
			second: &one,
			fail:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tt := new(testing.T)
			DeepEqual(tt, testCase.first, testCase.second)
			fail := tt.Failed()
			if testCase.fail != fail {
				t.Errorf("Expected failure: %v\nActual failure: %v", testCase.fail, fail)
			}
		})
	}
}
