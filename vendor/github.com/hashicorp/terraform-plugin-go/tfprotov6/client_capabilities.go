// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

// ConfigureProviderClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ConfigureProvider RPC,
// such as forward-compatible Terraform behavior changes.
type ConfigureProviderClientCapabilities struct {
	// DeferralAllowed signals that the request from Terraform is able to
	// handle deferred responses from the provider.
	DeferralAllowed bool
}

// ReadDataSourceClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ReadDataSource RPC,
// such as forward-compatible Terraform behavior changes.
type ReadDataSourceClientCapabilities struct {
	// DeferralAllowed signals that the request from Terraform is able to
	// handle deferred responses from the provider.
	DeferralAllowed bool
}

// ReadResourceClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ReadResource RPC,
// such as forward-compatible Terraform behavior changes.
type ReadResourceClientCapabilities struct {
	// DeferralAllowed signals that the request from Terraform is able to
	// handle deferred responses from the provider.
	DeferralAllowed bool
}

// PlanResourceChangeClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the PlanResourceChange RPC,
// such as forward-compatible Terraform behavior changes.
type PlanResourceChangeClientCapabilities struct {
	// DeferralAllowed signals that the request from Terraform is able to
	// handle deferred responses from the provider.
	DeferralAllowed bool
}

// ImportResourceStateClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ImportResourceState RPC,
// such as forward-compatible Terraform behavior changes.
type ImportResourceStateClientCapabilities struct {
	// DeferralAllowed signals that the request from Terraform is able to
	// handle deferred responses from the provider.
	DeferralAllowed bool
}
