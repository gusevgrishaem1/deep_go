package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.errs {
		buf.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	return fmt.Sprintf("%d errors occured:\n%s\n", len(e.errs), buf.String())
}

func Append(err error, errs ...error) *MultiError {
	var allErrs []error

	if err != nil {
		if mErr, ok := err.(*MultiError); ok {
			allErrs = append(allErrs, mErr.errs...)
		} else {
			allErrs = append(allErrs, err)
		}
	}

	allErrs = append(allErrs, errs...)

	return &MultiError{errs: allErrs}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func (e *MultiError) Unwrap() []error {
	return e.errs
}

// Тесты errors.Is и errors.As

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

func TestErrorsIsAndAs(t *testing.T) {
	err1 := &MyError{"error 1"}
	err2 := errors.New("error 2")

	var err error
	err = Append(err, err1)
	err = Append(err, err2)

	assert.True(t, errors.Is(err, err1), "errors.Is не нашёл err1")
	assert.True(t, errors.Is(err, err2), "errors.Is не нашёл err2")

	var target *MyError
	ok := errors.As(err, &target)
	assert.True(t, ok, "errors.As не нашёл MyError")
}
