// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import "time"

type Timeouts struct {
	DefaultCreateTimeout time.Duration
	DefaultReadTimeout   time.Duration
	DefaultUpdateTimeout time.Duration
	DefaultDeleteTimeout time.Duration
}
