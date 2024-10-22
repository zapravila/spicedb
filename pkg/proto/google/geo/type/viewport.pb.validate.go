// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: google/geo/type/viewport.proto

package _type

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Viewport with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Viewport) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Viewport with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ViewportMultiError, or nil
// if none found.
func (m *Viewport) ValidateAll() error {
	return m.validate(true)
}

func (m *Viewport) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetLow()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ViewportValidationError{
					field:  "Low",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ViewportValidationError{
					field:  "Low",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLow()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ViewportValidationError{
				field:  "Low",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetHigh()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ViewportValidationError{
					field:  "High",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ViewportValidationError{
					field:  "High",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetHigh()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ViewportValidationError{
				field:  "High",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ViewportMultiError(errors)
	}

	return nil
}

// ViewportMultiError is an error wrapping multiple validation errors returned
// by Viewport.ValidateAll() if the designated constraints aren't met.
type ViewportMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ViewportMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ViewportMultiError) AllErrors() []error { return m }

// ViewportValidationError is the validation error returned by
// Viewport.Validate if the designated constraints aren't met.
type ViewportValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ViewportValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ViewportValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ViewportValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ViewportValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ViewportValidationError) ErrorName() string { return "ViewportValidationError" }

// Error satisfies the builtin error interface
func (e ViewportValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sViewport.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ViewportValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ViewportValidationError{}