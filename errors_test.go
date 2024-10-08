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
			name:       "ErrorIsAll/Empty/Nil",
			errorCheck: ErrorIsAll(),
			err:        nil,
			fail:       false,
		},
		{
			name:       "ErrorIsAll/Empty/ErrTest",
			errorCheck: ErrorIsAll(),
			err:        ErrTest,
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
			tMock := newTMock()
			testCase.errorCheck(tMock, testCase.err)
			fail := tMock.ErrorfCalled > 0
			if testCase.fail != fail {
				t.Errorf("Expected failure: %v\nActual failure: %v", testCase.fail, fail)
			}
		})
	}
}

func TestMust(t *testing.T) {
	type testCase struct {
		f                  func() (bool, error)
		expectedValue      bool
		expectedFatalCalls int32
	}
	run := func(name string, testCase testCase) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			tMock := newTMock()
			value := Must(testCase.f())(tMock)
			Equal(t, value, testCase.expectedValue)
			Equal(t, testCase.expectedFatalCalls, tMock.FatalfCalled)
		})
	}

	run("Error", testCase{
		f: func() (bool, error) {
			return false, ErrTest
		},
		expectedValue:      false,
		expectedFatalCalls: 1,
	})
	run("Success", testCase{
		f: func() (bool, error) {
			return true, nil
		},
		expectedValue:      true,
		expectedFatalCalls: 0,
	})
}

func TestMust0(t *testing.T) {
	type testCase struct {
		f                  func() error
		expectedFatalCalls int32
	}
	run := func(name string, testCase testCase) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			tMock := newTMock()
			Must0(testCase.f())(tMock)
			Equal(t, testCase.expectedFatalCalls, tMock.FatalfCalled)
		})
	}

	run("Error", testCase{
		f: func() error {
			return ErrTest
		},
		expectedFatalCalls: 1,
	})
	run("Success", testCase{
		f: func() error {
			return nil
		},
		expectedFatalCalls: 0,
	})
}

func TestMust2(t *testing.T) {
	type testCase struct {
		f                  func() (bool, bool, error)
		expectedValue1     bool
		expectedValue2     bool
		expectedFatalCalls int32
	}
	run := func(name string, testCase testCase) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			tMock := newTMock()
			value1, value2 := Must2(testCase.f())(tMock)
			Equal(t, value1, testCase.expectedValue1)
			Equal(t, value2, testCase.expectedValue2)
			Equal(t, testCase.expectedFatalCalls, tMock.FatalfCalled)
		})
	}

	run("Error", testCase{
		f: func() (bool, bool, error) {
			return false, false, ErrTest
		},
		expectedValue1:     false,
		expectedValue2:     false,
		expectedFatalCalls: 1,
	})
	run("Success", testCase{
		f: func() (bool, bool, error) {
			return true, true, nil
		},
		expectedValue1:     true,
		expectedValue2:     true,
		expectedFatalCalls: 0,
	})
}

func TestMust3(t *testing.T) {
	type testCase struct {
		f                  func() (bool, bool, bool, error)
		expectedValue1     bool
		expectedValue2     bool
		expectedValue3     bool
		expectedFatalCalls int32
	}
	run := func(name string, testCase testCase) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			tMock := newTMock()
			value1, value2, value3 := Must3(testCase.f())(tMock)
			Equal(t, value1, testCase.expectedValue1)
			Equal(t, value2, testCase.expectedValue2)
			Equal(t, value3, testCase.expectedValue3)
			Equal(t, testCase.expectedFatalCalls, tMock.FatalfCalled)
		})
	}

	run("Error", testCase{
		f: func() (bool, bool, bool, error) {
			return false, false, false, ErrTest
		},
		expectedValue1:     false,
		expectedValue2:     false,
		expectedValue3:     false,
		expectedFatalCalls: 1,
	})
	run("Success", testCase{
		f: func() (bool, bool, bool, error) {
			return true, true, true, nil
		},
		expectedValue1:     true,
		expectedValue2:     true,
		expectedValue3:     true,
		expectedFatalCalls: 0,
	})
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

func newTMock() *TMock {
	return &TMock{
		HelperStub: func() {},
		ErrorfStub: func(format string, args ...any) {},
		FatalfStub: func(format string, args ...any) {},
	}
}
