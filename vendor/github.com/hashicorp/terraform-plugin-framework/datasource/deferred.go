// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

const (
	// DeferredReasonUnknown is used to indicate an invalid `DeferredReason`.
	// Provider developers should not use it.
	DeferredReasonUnknown DeferredReason = 0

	// DeferredReasonDataSourceConfigUnknown is used to indicate that the resource configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonDataSourceConfigUnknown DeferredReason = 1

	// DeferredReasonProviderConfigUnknown is used to indicate that the provider configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonProviderConfigUnknown DeferredReason = 2

	// DeferredReasonAbsentPrereq is used to indicate that a hard dependency has not been satisfied.
	DeferredReasonAbsentPrereq DeferredReason = 3
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
	case 1:
		return "Data Source Config Unknown"
	case 2:
		return "Provider Config Unknown"
	case 3:
		return "Absent Prerequisite"
	}
	return "Unknown"
}
