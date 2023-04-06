package optional

import (
	"fmt"
	"reflect"
	"testing"
)

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

	expect(t, !nothing.IsPresent(), "Nothing.IsPresent should return false")
	expect(t, nothing.Filter(False[string]) == nothing, "Nothing.Filter should return nothing")
	expect(t, nothing.OrElse("whatever") == "whatever", "Nothing.OrElse should return else value")
	_, err := nothing.Get()
	if err != ErrNoSuchElement {
		t.Errorf("Nothing.Get should return error")
	}
}

func TestJust(t *testing.T) {
	just := Just[string]{
		value: "whatever",
	}

	expect(t, just.IsPresent(), "Just.IsPresent should return true")
	expect(t, just.Filter(False[string]) == Nothing[string]{}, "Just.Filter should return nothing for false")
	expect(t, just.Filter(True[string]) == just, "Just.Filter should return itself for true")
	expect(t, just.OrElse("whatever2") == "whatever", "Just.OrElse should return its value")
	value, err := just.Get()
	if err == ErrNoSuchElement || value != "whatever" {
		t.Errorf("Nothing.Get should return error")
	}
}

func TestError(t *testing.T) {
	e := Error[string]{
		err: fmt.Errorf("Whatever"),
	}
	_, err := e.Get()
	if err.Error() != "Whatever" {
		t.Errorf("Error.Get should return internal error, but got %v", err)
	}
}
