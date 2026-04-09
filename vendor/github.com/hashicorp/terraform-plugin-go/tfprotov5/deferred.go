// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov5

const (
	// DeferredReasonUnknown is used to indicate an invalid `DeferredReason`.
	// Provider developers should not use it.
	DeferredReasonUnknown DeferredReason = 0

	// DeferredReasonResourceConfigUnknown is used to indicate that the resource configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonResourceConfigUnknown DeferredReason = 1

	// DeferredReasonProviderConfigUnknown is used to indicate that the provider configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonProviderConfigUnknown DeferredReason = 2

	// DeferredReasonAbsentPrereq is used to indicate that a hard dependency has not been satisfied.
	DeferredReasonAbsentPrereq DeferredReason = 3
)

// Deferred is used to indicate to Terraform that a change needs to be deferred for a reason.
type Deferred struct {
	// Reason is the reason for deferring the change.
	Reason DeferredReason
}

// DeferredReason represents different reasons for deferring a change.
type DeferredReason int32

func (d DeferredReason) String() string {
	switch d {
	case 0:
		return "UNKNOWN"
	case 1:
		return "RESOURCE_CONFIG_UNKNOWN"
	case 2:
		return "PROVIDER_CONFIG_UNKNOWN"
	case 3:
		return "ABSENT_PREREQ"
	}
	return "UNKNOWN"
}
