// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rules

type Rule interface {
	Run() []error
	Name() string
	Description() string
}
