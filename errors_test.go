package expect

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		name       string
		errorCheck ErrorCheck
		err        error
		fail       bool
	}{
		{
			name:       "ErrorNil/Nil",
			errorCheck: ErrorNil,
			err:        nil,
			fail:       false,
		},
		{
			name:       "ErrorNil/ErrTest",
			errorCheck: ErrorNil,
			err:        ErrTest,
			fail:       true,
		},
		{
			name:       "ErrorNonNil/Nil",
			errorCheck: ErrorNonNil,
			err:        nil,
			fail:       true,
		},
		{
			name:       "ErrorNonNil/ErrTest",
			errorCheck: ErrorNonNil,
			err:        ErrTest,
			fail:       false,
		},
		{
			name:       "ErrorIs/Nil/Nil",
			errorCheck: ErrorIs(nil),
			err:        nil,
			fail:       false,
		},
		{
			name:       "ErrorIs/Nil/ErrTest",
			errorCheck: ErrorIs(nil),
			err:        ErrTest,
			fail:       true,
		},
		{
			name:       "ErrorIs/ErrTest/Nil",
			errorCheck: ErrorIs(ErrTest),
			err:        nil,
			fail:       true,
		},
		{
			name:       "ErrorIs/ErrTest/ErrTest",
			errorCheck: ErrorIs(ErrTest),
			err:        ErrTest,
			fail:       false,
		},
		{
			name:       "ErrorIs/ErrTest/WrappedErrTest",
			errorCheck: ErrorIs(ErrTest),
			err:        fmt.Errorf("error: %w", ErrTest),
			fail:       false,
		},
		{
			name:       "ErrorIs/WrappedErrTest/ErrTest",
			errorCheck: ErrorIs(fmt.Errorf("error: %w", ErrTest)),
			err:        ErrTest,
			fail:       true,
		},
		{
			name:       "ErrorIs/ErrTest/Other",
			errorCheck: ErrorIs(ErrTest),
			err:        errors.New("other"),
			fail:       true,
		},
		{
			name:       "ErrorIs/Other/ErrTest",
			errorCheck: ErrorIs(errors.New("other")),
			err:        ErrTest,
			fail:       true,
		},
		{
			name:       "ErrorAs/TestErrorA/Nil",
			errorCheck: ErrorAs[testErrorA](),
			err:        nil,
			fail:       true,
		},
		{
			name:       "ErrorAs/TestErrorA/TestErrorA",
			errorCheck: ErrorAs[testErrorA](),
			err:        testErrorA{},
			fail:       false,
		},
		{
			name:       "ErrorAs/TestErrorA/TestErrorB",
			errorCheck: ErrorAs[testErrorA](),
			err:        testErrorB{},
			fail:       false,
		},
		{
			name:       "ErrorAs/TestErrorB/TestErrorA",
			errorCheck: ErrorAs[testErrorB](),
			err:        testErrorA{},
			fail:       true,
		},
		{
			name:       "ErrorIsAll/Nil/Nil",
			errorCheck: ErrorIsAll(nil),
			err:        nil,
			fail:       false,
		},
		{
			name:       "ErrorIsAll/Nil/ErrTest",
			errorCheck: ErrorIsAll(nil),
			err:        ErrTest,
			fail:       true,
		},
		{
			name:       "ErrorIsAll/ErrTest/Nil",
			errorCheck: ErrorIsAll(ErrTest),
			err:        nil,
			fail:       true,
		},
		{
			name:       "ErrorIsAll/ErrTest/ErrTest",
			errorCheck: ErrorIsAll(ErrTest),
			err:        ErrTest,
			fail:       false,
		},
		{
			name:       "ErrorIsAll/ErrTest/ErrTest/WrappedErrTest",
			errorCheck: ErrorIsAll(ErrTest),
			err:        errors.Join(ErrTest, fmt.Errorf("error: %w", ErrTest)),
			fail:       false,
		},
		{
			name:       "ErrorIsAll/ErrTest/ErrTest/Other",
			errorCheck: ErrorIsAll(ErrTest),
			err:        errors.Join(ErrTest, errors.New("other")),
			fail:       false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tt := new(testing.T)
			testCase.errorCheck(tt, testCase.err)
			fail := tt.Failed()
			if testCase.fail != fail {
				t.Errorf("Expected failure: %v\nActual failure: %v", testCase.fail, fail)
			}
		})
	}
}

type testErrorA struct{}

func (e testErrorA) Error() string {
	return "error A"
}

type testErrorB struct {
	a testErrorA
}

func (e testErrorB) Error() string {
	return "error B"
}

func (e testErrorB) As(target any) bool {
	switch t := target.(type) {
	case *testErrorA:
		*t = e.a
		return true
	default:
		return false
	}
}
