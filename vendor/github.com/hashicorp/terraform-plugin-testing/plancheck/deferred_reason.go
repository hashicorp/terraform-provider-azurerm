// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

// DeferredReason is a string stored in the plan file which indicates why Terraform
// is deferring a change for a resource.
type DeferredReason string

const (
	// DeferredReasonResourceConfigUnknown is used to indicate that the resource configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonResourceConfigUnknown DeferredReason = "resource_config_unknown"

	// DeferredReasonProviderConfigUnknown is used to indicate that the provider configuration
	// is partially unknown and the real values need to be known before the change can be planned.
	DeferredReasonProviderConfigUnknown DeferredReason = "provider_config_unknown"

	// DeferredReasonAbsentPrereq is used to indicate that a hard dependency has not been satisfied.
	DeferredReasonAbsentPrereq DeferredReason = "absent_prereq"
)
