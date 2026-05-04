// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

// MAINTAINER NOTE: The deferred reason at enum value 1 in the plugin-protocol
// is not relevant for provider-level automatic deferred responses.
// provider.DeferredReason is directly mapped to the plugin-protocol which is
// why enum value 1 is skipped here
const (
	// DeferredReasonUnknown is used to indicate an invalid `DeferredReason`.
	// Provider developers should not use it.
	DeferredReasonUnknown DeferredReason = 0

	// DeferredReasonProviderConfigUnknown is used to indicate that the provider configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonProviderConfigUnknown DeferredReason = 2
)

// Deferred is used to indicate to Terraform that a change needs to be deferred for a reason.
//
// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
// to change or break without warning. It is not protected by version compatibility guarantees.
type Deferred struct {
	// Reason is the reason for deferring the change.
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
