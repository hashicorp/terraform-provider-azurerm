// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

// anyFloat is a constraint that permits any floating-point type. This type
// definition is copied rather than depending on x/exp/constraints since the
// dependency is otherwise unneeded, the definition is relatively trivial and
// static, and the Go language maintainers are not sure if/where these will live
// in the standard library.
//
// Reference: https://github.com/golang/go/issues/61914
type anyFloat interface {
	~float32 | ~float64
}

// anyInteger is a constraint that permits any integer type. This type
// definition is copied rather than depending on x/exp/constraints since the
// dependency is otherwise unneeded, the definition is relatively trivial and
// static, and the Go language maintainers are not sure if/where these will live
// in the standard library.
//
// Reference: https://github.com/golang/go/issues/61914
type anyInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
