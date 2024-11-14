// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

// MAINTAINER NOTE: Only PROVIDER_CONFIG_UNKNOWN (enum value 2 in the plugin-protocol) is relevant
// for SDKv2. Since (Deferred).Reason is mapped directly to the plugin-protocol,
// the other enum values are intentionally omitted here.
const (
	// DeferredReasonUnknown is used to indicate an invalid `DeferredReason`.
	// Provider developers should not use it.
	DeferredReasonUnknown DeferredReason = 0

	// DeferredReasonProviderConfigUnknown represents a deferred reason caused
	// by unknown provider configuration.
	DeferredReasonProviderConfigUnknown DeferredReason = 2
)

// Deferred is used to indicate to Terraform that a resource or data source is not able
// to be applied yet and should be skipped (deferred). After completing an apply that has deferred actions,
// the practitioner can then execute additional plan and apply “rounds” to eventually reach convergence
// where there are no remaining deferred actions.
//
// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
// to change or break without warning. It is not protected by version compatibility guarantees.
type Deferred struct {
	// Reason represents the deferred reason.
	Reason DeferredReason
}

// DeferredReason represents different reasons for deferring a change.
//
// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
// to change or break without warning. It is not protected by version compatibility guarantees.
type DeferredReason int32

func (d DeferredReason) String() string {
	switch d {
	case 0:
		return "Unknown"
	case 2:
		return "Provider Config Unknown"
	}
	return "Unknown"
}
