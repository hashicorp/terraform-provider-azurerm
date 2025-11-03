// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

// ValidateSchemaClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the schema validation
// RPCs, such as forward-compatible Terraform behavior changes.
type ValidateSchemaClientCapabilities struct {
	// WriteOnlyAttributesAllowed indicates that the Terraform client
	// initiating the request supports write-only attributes for managed
	// resources.
	//
	// This client capability is only populated during managed resource schema
	// validation.
	WriteOnlyAttributesAllowed bool
}
