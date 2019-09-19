// Code generated by protoc-gen-validate
// source: envoy/api/v2/core/address.proto
// DO NOT EDIT!!!

package core

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
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
	_ = types.DynamicAny{}
)

// Validate checks the field values on Pipe with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Pipe) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetPath()) < 1 {
		return PipeValidationError{
			Field:  "Path",
			Reason: "value length must be at least 1 bytes",
		}
	}

	return nil
}

// PipeValidationError is the validation error returned by Pipe.Validate if the
// designated constraints aren't met.
type PipeValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e PipeValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPipe.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = PipeValidationError{}

// Validate checks the field values on SocketAddress with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *SocketAddress) Validate() error {
	if m == nil {
		return nil
	}

	if _, ok := SocketAddress_Protocol_name[int32(m.GetProtocol())]; !ok {
		return SocketAddressValidationError{
			Field:  "Protocol",
			Reason: "value must be one of the defined enum values",
		}
	}

	if len(m.GetAddress()) < 1 {
		return SocketAddressValidationError{
			Field:  "Address",
			Reason: "value length must be at least 1 bytes",
		}
	}

	// no validation rules for ResolverName

	// no validation rules for Ipv4Compat

	switch m.PortSpecifier.(type) {

	case *SocketAddress_PortValue:

		if m.GetPortValue() > 65535 {
			return SocketAddressValidationError{
				Field:  "PortValue",
				Reason: "value must be less than or equal to 65535",
			}
		}

	case *SocketAddress_NamedPort:
		// no validation rules for NamedPort

	default:
		return SocketAddressValidationError{
			Field:  "PortSpecifier",
			Reason: "value is required",
		}

	}

	return nil
}

// SocketAddressValidationError is the validation error returned by
// SocketAddress.Validate if the designated constraints aren't met.
type SocketAddressValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e SocketAddressValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketAddress.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = SocketAddressValidationError{}

// Validate checks the field values on TcpKeepalive with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *TcpKeepalive) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetKeepaliveProbes()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TcpKeepaliveValidationError{
				Field:  "KeepaliveProbes",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetKeepaliveTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TcpKeepaliveValidationError{
				Field:  "KeepaliveTime",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetKeepaliveInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TcpKeepaliveValidationError{
				Field:  "KeepaliveInterval",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// TcpKeepaliveValidationError is the validation error returned by
// TcpKeepalive.Validate if the designated constraints aren't met.
type TcpKeepaliveValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e TcpKeepaliveValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTcpKeepalive.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = TcpKeepaliveValidationError{}

// Validate checks the field values on BindConfig with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *BindConfig) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetSourceAddress()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BindConfigValidationError{
				Field:  "SourceAddress",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetFreebind()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BindConfigValidationError{
				Field:  "Freebind",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	for idx, item := range m.GetSocketOptions() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return BindConfigValidationError{
					Field:  fmt.Sprintf("SocketOptions[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// BindConfigValidationError is the validation error returned by
// BindConfig.Validate if the designated constraints aren't met.
type BindConfigValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e BindConfigValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBindConfig.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = BindConfigValidationError{}

// Validate checks the field values on Address with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Address) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Address.(type) {

	case *Address_SocketAddress:

		if v, ok := interface{}(m.GetSocketAddress()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AddressValidationError{
					Field:  "SocketAddress",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *Address_Pipe:

		if v, ok := interface{}(m.GetPipe()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AddressValidationError{
					Field:  "Pipe",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	default:
		return AddressValidationError{
			Field:  "Address",
			Reason: "value is required",
		}

	}

	return nil
}

// AddressValidationError is the validation error returned by Address.Validate
// if the designated constraints aren't met.
type AddressValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e AddressValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddress.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = AddressValidationError{}

// Validate checks the field values on CidrRange with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *CidrRange) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetAddressPrefix()) < 1 {
		return CidrRangeValidationError{
			Field:  "AddressPrefix",
			Reason: "value length must be at least 1 bytes",
		}
	}

	if wrapper := m.GetPrefixLen(); wrapper != nil {

		if wrapper.GetValue() > 128 {
			return CidrRangeValidationError{
				Field:  "PrefixLen",
				Reason: "value must be less than or equal to 128",
			}
		}

	}

	return nil
}

// CidrRangeValidationError is the validation error returned by
// CidrRange.Validate if the designated constraints aren't met.
type CidrRangeValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e CidrRangeValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCidrRange.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = CidrRangeValidationError{}
