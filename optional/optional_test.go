package optional

import (
	"fmt"
	"strconv"
	"strings"

	"reflect"
	"testing"
)

func doNothing() {}

func doNothingStr(string) {}

func doNothingErr(error) {}

type someNullable struct {
	someValue string
}

type someOtherNullable struct {
	someValue int
}

func expect(t *testing.T, value bool, errmsg string) {
	if !value {
		t.Error(errmsg)
	}
}

func errtuple[T any](value T, err error) (T, error) {
	return value, err
}

func TestOf(t *testing.T) {

	t.Run("OfNullable(nil) returns Nothing", func(t *testing.T) {
		var nilt *string = nil

		if !reflect.DeepEqual(OfNullable(nilt), Nothing[*string]{}) {
			t.Errorf("OfNullable[...](%v) expected to return Nothing", nilt)
		}
	})

	t.Run("OfNullable(value) returns Just", func(t *testing.T) {
		nilt := "Whatever"
		expected := "Whatever"

		if !reflect.DeepEqual(OfNullable(&nilt), Just[*string]{value: &expected}) {
			t.Errorf("OfNullable[...](%v) expected to return Just", nilt)
		}
	})

	t.Run("Of (Value) returns Just", func(t *testing.T) {
		if !reflect.DeepEqual(Of("Whatever"), Just[string]{value: "Whatever"}) {
			t.Errorf("Of[...](%v) expected to return Just", "Whatever")
		}
	})

	t.Run("OfError returns Just if error is nil", func(t *testing.T) {

		if !reflect.DeepEqual(OfError(errtuple("A", nil)), Just[string]{value: "A"}) {
			t.Error("OfError[...]('A', nil) expected to return Just")
		}
	})

	t.Run("OfError returns Error if error is not nil", func(t *testing.T) {

		if !reflect.DeepEqual(OfError(errtuple("A", fmt.Errorf("B"))), Error[string]{err: fmt.Errorf("B")}) {
			t.Error("OfError[...]('A', err) expected to return Error")
		}
	})

	t.Run("OfErrorNullable returns Just if error is nil", func(t *testing.T) {
		val := "Whatever"
		if !reflect.DeepEqual(OfErrorNullable(errtuple(&val, nil)), Just[*string]{value: &val}) {
			t.Error("OfErrorNullable[...]('A', nil) expected to return Just")
		}
	})

	t.Run("OfErrorNullable returns Error if error is not nil", func(t *testing.T) {
		val := "Whatever"
		if !reflect.DeepEqual(OfErrorNullable(errtuple(&val, fmt.Errorf("B"))), Error[*string]{err: fmt.Errorf("B")}) {
			t.Error("OfErrorNullable[...]('A', err) expected to return Error")
		}
	})
}

func TestNothing(t *testing.T) {
	nothing := Nothing[string]{}

	t.Run("Nothing basic methods implementation", func(t *testing.T) {
		expect(t, !nothing.IsPresent(), "Nothing.IsPresent should return false")
		expect(t, nothing.Filter(False[string]) == nothing, "Nothing.Filter should return nothing")
		expect(t, nothing.OrElse("whatever") == "whatever", "Nothing.OrElse should return else value")
		expect(t, nothing.Map(strings.ToLower) == nothing, "Nothing.Map should return nothing")
		expect(t, nothing.GetError() == nil, "Nothing.GetError should return nil")
		expect(t, nothing.IfPresent(doNothingStr) == nothing, "Nothing.IfPresent(...) returns itself")
		expect(t, nothing.IfError(doNothingErr) == nothing, "Nothing.IfError(...) returns itself")
		expect(t, nothing.IfNothing(doNothing) == nothing, "Nothing.IfNothing(...) returns itself")
	})

	t.Run("Nothing.Get returns NoSuchElement error", func(t *testing.T) {
		_, err := nothing.Get()
		if err != ErrNoSuchElement {
			t.Errorf("Nothing.Get should return error")
		}
	})

	t.Run("Nothing.IfPresent does not call function", func(t *testing.T) {
		ptr := new(string)
		*ptr = "123456"
		nothing.IfPresent(func(s string) {
			*ptr = "ifpresent"
		})
		expect(t, *ptr == "123456", "Nothing.IfPresent should not be invoked")
	})

	t.Run("Nothing.If invokes only onEmpty function", func(t *testing.T) {
		invoked := false
		nothing.If(
			func(s string) {
				t.Errorf("Nothing not expected to invoke If-present: %s", s)
			},
			func(e error) {
				t.Errorf("Nothing not expected to invoke If-error: %s", e)
			},
			func() {
				invoked = true
			})
		if !invoked {
			t.Errorf("Nothing expected to invoke If-nothing")
		}
	})

	t.Run("Nothing.IfNothing invokes function", func(t *testing.T) {
		invoked := false
		nothing.IfNothing(func() {
			invoked = true
		})
		if !invoked {
			t.Errorf("Nothing expected to invoke If-nothing")
		}
	})

	t.Run("Nothing.IfError does not invoke function", func(t *testing.T) {
		ptr := new(error)
		*ptr = fmt.Errorf("An error")
		nothing.IfError(func(e error) {
			*ptr = e
		})
		expect(t, (*ptr).Error() == "An error", "Nothing.IfError should not be invoked")
	})
}

func TestJust(t *testing.T) {
	just := Just[string]{
		value: "whatever",
	}

	t.Run("Just basic methods implementation", func(t *testing.T) {
		expect(t, just.IsPresent(), "Just.IsPresent should return true")
		expect(t, just.Filter(False[string]) == Nothing[string]{}, "Just.Filter should return nothing for false")
		expect(t, just.Filter(True[string]) == just, "Just.Filter should return itself for true")
		expect(t, just.OrElse("whatever2") == "whatever", "Just.OrElse should return its value")
		expect(t, just.GetError() == nil, "Just.GetError should return nil")
		expect(t, just.IfPresent(doNothingStr) == just, "Just.IfPresent(...) returns itself")
		expect(t, just.IfError(doNothingErr) == just, "Just.IfError(...) returns itself")
		expect(t, just.IfNothing(doNothing) == just, "Just.IfNothing(...) returns itself")
	})

	t.Run("Just.Map returns updated version", func(t *testing.T) {
		upper, _ := just.Map(strings.ToUpper).Get()
		if upper != "WHATEVER" {
			t.Errorf("Just.Map expected to change value")
		}
	})

	t.Run("Just.Get returns value and nil for error", func(t *testing.T) {
		value, err := just.Get()
		if err == ErrNoSuchElement || value != "whatever" {
			t.Errorf("Just.Get should return value and nil error")
		}
	})

	t.Run("Just.IfPresent invokes function", func(t *testing.T) {
		ptr := new(string)
		*ptr = "123456"
		just.IfPresent(func(s string) {
			*ptr = s
		})
		expect(t, *ptr == just.value, "Just.IfPresent should be invoked")
	})

	t.Run("Just.If invokes only IfPresent function", func(t *testing.T) {
		ptr := new(string)
		*ptr = "123456"
		just.If(
			func(s string) {
				*ptr = s
			},
			func(e error) {
				t.Errorf("Just not expected to invoke If-error: %s", e)
			},
			func() {
				t.Errorf("Just not expected to invoke If-nothing")
			})
		expect(t, *ptr == just.value, "Just.If-present should be invoked")
	})

	t.Run("Just.IfNothing does not invoke function", func(t *testing.T) {
		invoked := false
		just.IfNothing(func() {
			invoked = true
		})
		if invoked {
			t.Errorf("Just is not expected to invoke If-nothing")
		}
	})

	t.Run("Just.IfError does not invoke function", func(t *testing.T) {
		ptr := new(error)
		*ptr = fmt.Errorf("An error")
		just.IfError(func(e error) {
			*ptr = e
		})
		expect(t, (*ptr).Error() == "An error", "Just.IfPresent should not be invoked")
	})

}

func TestError(t *testing.T) {
	e := Error[string]{
		err: fmt.Errorf("Whatever"),
	}

	t.Run("Error basic methods implementation", func(t *testing.T) {
		expect(t, e.GetError().Error() == "Whatever", "Error.GetError should return error")
		expect(t, e.IfPresent(doNothingStr) == e, "Error.IfPresent(...) returns itself")
		expect(t, e.IfError(doNothingErr) == e, "Error.IfError(...) returns itself")
		expect(t, e.IfNothing(doNothing) == e, "Error.IfNothing(...) returns itself")
	})

	t.Run("Error.Get() returns error value", func(t *testing.T) {
		_, err := e.Get()
		if err.Error() != "Whatever" {
			t.Errorf("Error.Get should return internal error, but got %v", err)
		}
	})

	t.Run("Error.If invokes only onError function", func(t *testing.T) {
		ptr := new(error)
		*ptr = fmt.Errorf("An error")
		e.If(
			func(s string) {
				t.Errorf("Error not expected to invoke If-present: %s", s)
			},
			func(e error) {
				*ptr = e
			},
			func() {
				t.Errorf("Error not expected to invoke If-nothing")
			})
		expect(t, *ptr == e.err, "Just.If-error should be invoked")
	})

	t.Run("Error.IfNothing does not invoke function", func(t *testing.T) {
		invoked := false
		e.IfNothing(func() {
			invoked = true
		})
		if invoked {
			t.Errorf("Just is not expected to invoke If-nothing")
		}
	})

	t.Run("Error.IfError does invokes function", func(t *testing.T) {
		ptr := new(error)
		*ptr = fmt.Errorf("An error")
		e.IfError(func(e error) {
			*ptr = e
		})
		expect(t, (*ptr).Error() == e.err.Error(), "Error.IfError should be invoked")
	})
}

func TestMap(t *testing.T) {

	for _, tc := range []struct {
		desc   string
		opt    Optional[string]
		mapper func(string) int
		want   Optional[int]
	}{
		{
			desc: "Map of present",
			opt:  Of("123"),
			mapper: func(get string) int {
				i, _ := strconv.Atoi(get)
				return i
			},
			want: Of(123),
		},
		{
			desc: "Map of nothing returns nothging of different type",
			opt:  Nothing[string]{},
			mapper: func(get string) int {
				i, _ := strconv.Atoi(get)
				return i
			},
			want: Nothing[int]{},
		},
		{
			desc: "Map of Error returns Error of different type",
			opt:  Error[string]{err: fmt.Errorf("Error message")},
			mapper: func(get string) int {
				i, _ := strconv.Atoi(get)
				return i
			},
			want: Error[int]{err: fmt.Errorf("Error message")},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			result := Map(tc.opt, tc.mapper)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Map(%v, mapper) is expected to be %v, but got %v", tc.opt, tc.want, result)
			}
		})
	}

	for _, tc := range []struct {
		desc   string
		opt    Optional[*someNullable]
		mapper func(*someNullable) *someOtherNullable
		want   Optional[*someOtherNullable]
	}{
		{
			desc: "Map of nullable present",
			opt:  OfNullable(&someNullable{someValue: "123"}),
			mapper: func(get *someNullable) *someOtherNullable {
				i, _ := strconv.Atoi(get.someValue)
				return &someOtherNullable{someValue: i}
			},
			want: OfNullable(&someOtherNullable{someValue: 123}),
		},
		{
			desc: "Map of nullable nothing returns nothging of different type",
			opt:  Nothing[*someNullable]{},
			mapper: func(get *someNullable) *someOtherNullable {
				i, _ := strconv.Atoi(get.someValue)
				return &someOtherNullable{someValue: i}
			},
			want: Nothing[*someOtherNullable]{},
		},
		{
			desc: "Map of nullable Error returns Error of different type",
			opt:  Error[*someNullable]{err: fmt.Errorf("Error message")},
			mapper: func(get *someNullable) *someOtherNullable {
				i, _ := strconv.Atoi(get.someValue)
				return &someOtherNullable{someValue: i}
			},
			want: Error[*someOtherNullable]{err: fmt.Errorf("Error message")},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			result := Map(tc.opt, tc.mapper)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Map(%v, mapper) is expected to be %v, but got %v", tc.opt, tc.want, result)
			}
		})
	}

}
