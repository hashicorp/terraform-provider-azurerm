// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attr

import "fmt"

const (
	// ValueStateNull represents a value which is null.
	//
	// This value is 0 so it is the zero-value for types implementations.
	ValueStateNull ValueState = 0

	// ValueStateUnknown represents a value which is unknown.
	ValueStateUnknown ValueState = 1

	// ValueStateKnown represents a value which is known (not null or unknown).
	ValueStateKnown ValueState = 2
)

type ValueState uint8

func (s ValueState) String() string {
	switch s {
	case ValueStateKnown:
		return "known"
	case ValueStateNull:
		return "null"
	case ValueStateUnknown:
		return "unknown"
	default:
		panic(fmt.Sprintf("unhandled ValueState in String: %d", s))
	}
}
