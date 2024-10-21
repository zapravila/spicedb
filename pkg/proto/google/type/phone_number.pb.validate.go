// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: google/type/phone_number.proto

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

// Validate checks the field values on PhoneNumber with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PhoneNumber) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PhoneNumber with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PhoneNumberMultiError, or
// nil if none found.
func (m *PhoneNumber) ValidateAll() error {
	return m.validate(true)
}

func (m *PhoneNumber) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Extension

	switch v := m.Kind.(type) {
	case *PhoneNumber_E164Number:
		if v == nil {
			err := PhoneNumberValidationError{
				field:  "Kind",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		// no validation rules for E164Number
	case *PhoneNumber_ShortCode_:
		if v == nil {
			err := PhoneNumberValidationError{
				field:  "Kind",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetShortCode()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PhoneNumberValidationError{
						field:  "ShortCode",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PhoneNumberValidationError{
						field:  "ShortCode",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetShortCode()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PhoneNumberValidationError{
					field:  "ShortCode",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return PhoneNumberMultiError(errors)
	}

	return nil
}

// PhoneNumberMultiError is an error wrapping multiple validation errors
// returned by PhoneNumber.ValidateAll() if the designated constraints aren't met.
type PhoneNumberMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PhoneNumberMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PhoneNumberMultiError) AllErrors() []error { return m }

// PhoneNumberValidationError is the validation error returned by
// PhoneNumber.Validate if the designated constraints aren't met.
type PhoneNumberValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PhoneNumberValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PhoneNumberValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PhoneNumberValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PhoneNumberValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PhoneNumberValidationError) ErrorName() string { return "PhoneNumberValidationError" }

// Error satisfies the builtin error interface
func (e PhoneNumberValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPhoneNumber.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PhoneNumberValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PhoneNumberValidationError{}

// Validate checks the field values on PhoneNumber_ShortCode with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PhoneNumber_ShortCode) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PhoneNumber_ShortCode with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PhoneNumber_ShortCodeMultiError, or nil if none found.
func (m *PhoneNumber_ShortCode) ValidateAll() error {
	return m.validate(true)
}

func (m *PhoneNumber_ShortCode) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RegionCode

	// no validation rules for Number

	if len(errors) > 0 {
		return PhoneNumber_ShortCodeMultiError(errors)
	}

	return nil
}

// PhoneNumber_ShortCodeMultiError is an error wrapping multiple validation
// errors returned by PhoneNumber_ShortCode.ValidateAll() if the designated
// constraints aren't met.
type PhoneNumber_ShortCodeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PhoneNumber_ShortCodeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PhoneNumber_ShortCodeMultiError) AllErrors() []error { return m }

// PhoneNumber_ShortCodeValidationError is the validation error returned by
// PhoneNumber_ShortCode.Validate if the designated constraints aren't met.
type PhoneNumber_ShortCodeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PhoneNumber_ShortCodeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PhoneNumber_ShortCodeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PhoneNumber_ShortCodeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PhoneNumber_ShortCodeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PhoneNumber_ShortCodeValidationError) ErrorName() string {
	return "PhoneNumber_ShortCodeValidationError"
}

// Error satisfies the builtin error interface
func (e PhoneNumber_ShortCodeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPhoneNumber_ShortCode.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PhoneNumber_ShortCodeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PhoneNumber_ShortCodeValidationError{}
