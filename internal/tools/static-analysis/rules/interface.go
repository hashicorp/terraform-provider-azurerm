// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

type Rule interface {
	Run() []error
	Name() string
	Description() string
}
