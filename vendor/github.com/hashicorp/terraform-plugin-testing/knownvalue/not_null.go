// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = notNull{}

type notNull struct{}

// CheckValue determines whether the passed value is nil.
func (v notNull) CheckValue(other any) error {
	if other == nil {
		return fmt.Errorf("expected non-nil value for NotNull check, got: %T", other)
	}

	return nil
}

// String returns the string representation of notNull.
func (v notNull) String() string {
	return "not-null"
}

// NotNull returns a Check for asserting the value passed
// to the CheckValue method is not nil.
func NotNull() notNull {
	return notNull{}
}
