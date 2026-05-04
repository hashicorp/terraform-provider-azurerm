// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

// ServerCapabilities is a combination of tfprotov5.ServerCapabilities and
// tfprotov6.ServerCapabilties, which may diverge over time. If that happens,
// the toproto5 conversion logic will handle the appropriate filtering and the
// proto5server/fwserver logic will need to account for missing features.
type ServerCapabilities struct {
	// GetProviderSchemaOptional signals that the provider does not require the
	// GetProviderSchema RPC before other RPCs.
	//
	// This should always be enabled in framework providers and requires
	// Terraform 1.6 or later.
	GetProviderSchemaOptional bool

	// MoveResourceState signals that the provider is ready for the
	// MoveResourceState RPC.
	//
	// This should always be enabled in framework providers and requires
	// Terraform 1.8 or later.
	MoveResourceState bool

	// PlanDestroy signals that the provider is ready for the
	// PlanResourceChange RPC on resource destruction.
	//
	// This should always be enabled in framework providers and requires
	// Terraform 1.3 or later.
	PlanDestroy bool
}

// ServerCapabilities returns the server capabilities.
func (s *Server) ServerCapabilities() *ServerCapabilities {
	return &ServerCapabilities{
		GetProviderSchemaOptional: true,
		MoveResourceState:         true,
		PlanDestroy:               true,
	}
}
