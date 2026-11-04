// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

var _ statestore.StateStore = &StateStore{}

type StateStore struct {
	configuredChunkSize int64

	SchemaResponse         *statestore.SchemaResponse
	ConfigureResponse      *statestore.ConfigureResponse
	ValidateConfigResponse *statestore.ValidateConfigResponse

	// Declaring a mock state store is slightly more complicated then other types since the implementation is likely
	// needed to be stateful to work with multiple terraform commands.
	GetStatesFunc       func(context.Context, statestore.GetStatesRequest, *statestore.GetStatesResponse)
	DeleteStateFunc     func(context.Context, statestore.DeleteStateRequest, *statestore.DeleteStateResponse)
	LockStateFunc       func(context.Context, statestore.LockStateRequest, *statestore.LockStateResponse)
	UnlockStateFunc     func(context.Context, statestore.UnlockStateRequest, *statestore.UnlockStateResponse)
	ReadStateBytesFunc  func(context.Context, statestore.ReadStateBytesRequest, *statestore.ReadStateBytesResponse)
	WriteStateBytesFunc func(context.Context, statestore.WriteStateBytesRequest, *statestore.WriteStateBytesResponse)
}

func (s *StateStore) Schema(ctx context.Context, req statestore.SchemaRequest, resp *statestore.SchemaResponse) {
	if s.SchemaResponse != nil {
		resp.Diagnostics = s.SchemaResponse.Diagnostics
		resp.Schema = s.SchemaResponse.Schema
	}
}

func (s *StateStore) Configure(ctx context.Context, req statestore.ConfigureRequest, resp *statestore.ConfigureResponse) {
	if s.ConfigureResponse != nil {
		resp.Diagnostics = s.ConfigureResponse.Diagnostics

		if s.ConfigureResponse.ServerCapabilities != nil {
			resp.ServerCapabilities = s.ConfigureResponse.ServerCapabilities
		}
	}

	// Store configured chunk size
	s.configuredChunkSize = resp.ServerCapabilities.ChunkSize
}

func (s *StateStore) ConfiguredChunkSize() int64 {
	return s.configuredChunkSize
}

func (s *StateStore) ValidateConfig(ctx context.Context, req statestore.ValidateConfigRequest, resp *statestore.ValidateConfigResponse) {
	if s.ValidateConfigResponse != nil {
		resp.Diagnostics = s.ValidateConfigResponse.Diagnostics
	}
}

func (s *StateStore) GetStates(ctx context.Context, req statestore.GetStatesRequest, resp *statestore.GetStatesResponse) {
	if s.GetStatesFunc != nil {
		s.GetStatesFunc(ctx, req, resp)
	}
}

func (s *StateStore) DeleteState(ctx context.Context, req statestore.DeleteStateRequest, resp *statestore.DeleteStateResponse) {
	if s.DeleteStateFunc != nil {
		s.DeleteStateFunc(ctx, req, resp)
	}
}

func (s *StateStore) LockState(ctx context.Context, req statestore.LockStateRequest, resp *statestore.LockStateResponse) {
	if s.LockStateFunc != nil {
		s.LockStateFunc(ctx, req, resp)
	}
}

func (s *StateStore) UnlockState(ctx context.Context, req statestore.UnlockStateRequest, resp *statestore.UnlockStateResponse) {
	if s.UnlockStateFunc != nil {
		s.UnlockStateFunc(ctx, req, resp)
	}
}

func (s *StateStore) ReadStateBytes(ctx context.Context, req statestore.ReadStateBytesRequest, resp *statestore.ReadStateBytesResponse) {
	if s.ReadStateBytesFunc != nil {
		s.ReadStateBytesFunc(ctx, req, resp)
	}
}

func (s *StateStore) WriteStateBytes(ctx context.Context, req statestore.WriteStateBytesRequest, resp *statestore.WriteStateBytesResponse) {
	if s.WriteStateBytesFunc != nil {
		s.WriteStateBytesFunc(ctx, req, resp)
	}
}
