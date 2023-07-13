// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package state

import "strings"

// IgnoreCase is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func IgnoreCase(val interface{}) string {
	return strings.ToLower(val.(string))
}
