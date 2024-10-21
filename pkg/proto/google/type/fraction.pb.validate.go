// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: google/type/fraction.proto

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

// Validate checks the field values on Fraction with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Fraction) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Fraction with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in FractionMultiError, or nil
// if none found.
func (m *Fraction) ValidateAll() error {
	return m.validate(true)
}

func (m *Fraction) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Numerator

	// no validation rules for Denominator

	if len(errors) > 0 {
		return FractionMultiError(errors)
	}

	return nil
}

// FractionMultiError is an error wrapping multiple validation errors returned
// by Fraction.ValidateAll() if the designated constraints aren't met.
type FractionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FractionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FractionMultiError) AllErrors() []error { return m }

// FractionValidationError is the validation error returned by
// Fraction.Validate if the designated constraints aren't met.
type FractionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FractionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FractionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FractionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FractionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FractionValidationError) ErrorName() string { return "FractionValidationError" }

// Error satisfies the builtin error interface
func (e FractionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFraction.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FractionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FractionValidationError{}
