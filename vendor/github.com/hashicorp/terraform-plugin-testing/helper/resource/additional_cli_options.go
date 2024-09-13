// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

// AdditionalCLIOptions allows an intentionally limited set of options to be passed
// to the Terraform CLI when executing test steps.
type AdditionalCLIOptions struct {
	// Apply represents options to be passed to the `terraform apply` command.
	Apply ApplyOptions

	// Plan represents options to be passed to the `terraform plan` command.
	Plan PlanOptions
}

// ApplyOptions represents options to be passed to the `terraform apply` command.
type ApplyOptions struct {
	// AllowDeferral will pass the experimental `-allow-deferral` flag to the apply command.
	AllowDeferral bool
}

// PlanOptions represents options to be passed to the `terraform plan` command.
type PlanOptions struct {
	// AllowDeferral will pass the experimental `-allow-deferral` flag to the plan command.
	AllowDeferral bool
}
