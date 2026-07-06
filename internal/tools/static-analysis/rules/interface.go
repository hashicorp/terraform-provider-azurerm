// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

type Rule interface {
	Run() []error
	Name() string
	Description() string
}
