// Copyright IBM Corp. 2021, 2026
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

	// GenerateResourceConfig signals that the provider is ready for the
	// GenerateResourceConfig RPC.
	//
	// This should always be enabled in framework providers and requires
	// Terraform 1.14 or later.
	GenerateResourceConfig bool
}

// ServerCapabilities returns the server capabilities.
func (s *Server) ServerCapabilities() *ServerCapabilities {
	return &ServerCapabilities{
		GetProviderSchemaOptional: true,
		MoveResourceState:         true,
		PlanDestroy:               true,
		GenerateResourceConfig:    true,
	}
}

// StateStoreServerCapabilities is internal to fwserver as we don't need to expose it to state store implementations currently.
type StateStoreServerCapabilities struct {
	// ChunkSize is the provider-chosen size of state byte chunks that will be sent between Terraform and
	// the provider in the ReadStateBytes and WriteStateBytes RPC calls.
	//
	// As we don't expose this to providers during ConfigureStateStore currently, the provider-chosen size will always be
	// the Terraform core defaulted value (8 MB).
	ChunkSize int64
}
