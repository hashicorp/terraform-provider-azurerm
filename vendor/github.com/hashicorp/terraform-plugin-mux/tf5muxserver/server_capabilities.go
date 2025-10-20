// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import "github.com/hashicorp/terraform-plugin-go/tfprotov5"

// serverCapabilities always announces all ServerCapabilities. Individual
// capabilities are handled in their respective RPCs to protect downstream
// servers if they are not compatible with a capability.
var serverCapabilities = &tfprotov5.ServerCapabilities{
	GetProviderSchemaOptional: true,
	MoveResourceState:         true,
	PlanDestroy:               true,
}

// serverSupportsPlanDestroy returns true if the given ServerCapabilities is not
// nil and enables the PlanDestroy capability.
func serverSupportsPlanDestroy(capabilities *tfprotov5.ServerCapabilities) bool {
	if capabilities == nil {
		return false
	}

	return capabilities.PlanDestroy
}
