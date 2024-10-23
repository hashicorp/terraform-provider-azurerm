// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

const (
	// DataDescriptionConfiguration is used for Data that represents
	// a configuration-based value.
	DataDescriptionConfiguration DataDescription = "configuration"

	// DataDescriptionPlan is used for Data that represents
	// a plan-based value.
	DataDescriptionPlan DataDescription = "plan"

	// DataDescriptionState is used for Data that represents
	// a state-based value.
	DataDescriptionState DataDescription = "state"
)

// DataDescription is a human friendly type for Data. Used in error
// diagnostics.
type DataDescription string

// String returns the lowercase string of the description.
func (d DataDescription) String() string {
	switch d {
	case "":
		return "data"
	default:
		return string(d)
	}
}

// Title returns the titlecase string of the description.
func (d DataDescription) Title() string {
	switch d {
	case DataDescriptionConfiguration:
		return "Configuration"
	case DataDescriptionPlan:
		return "Plan"
	case DataDescriptionState:
		return "State"
	default:
		return "Data"
	}
}
