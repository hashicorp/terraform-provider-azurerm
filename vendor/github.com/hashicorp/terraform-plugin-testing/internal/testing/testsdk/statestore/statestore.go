// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type StateStore interface {
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
	GetStates(context.Context, GetStatesRequest, *GetStatesResponse)
	DeleteState(context.Context, DeleteStateRequest, *DeleteStateResponse)
	LockState(context.Context, LockStateRequest, *LockStateResponse)
	UnlockState(context.Context, UnlockStateRequest, *UnlockStateResponse)

	// For ease-of-use, the streaming of chunk data is handled in the provider server for reading/writing
	ReadStateBytes(context.Context, ReadStateBytesRequest, *ReadStateBytesResponse)
	WriteStateBytes(context.Context, WriteStateBytesRequest, *WriteStateBytesResponse)

	// This isn't a GRPC call, but it allows the implementation to define the chunk size while the provider server owns the actual chunking logic.
	ConfiguredChunkSize() int64
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}

type ConfigureRequest struct {
	Config             tftypes.Value
	ClientCapabilities *tfprotov6.ConfigureStateStoreClientCapabilities
}

type ConfigureResponse struct {
	Diagnostics        []*tfprotov6.Diagnostic
	ServerCapabilities *tfprotov6.StateStoreServerCapabilities
}

type ValidateConfigRequest struct {
	Config tftypes.Value
}

type ValidateConfigResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type GetStatesRequest struct{}

type GetStatesResponse struct {
	StateIDs    []string
	Diagnostics []*tfprotov6.Diagnostic
}

type DeleteStateRequest struct {
	StateID string
}

type DeleteStateResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type LockStateRequest struct {
	StateID   string
	Operation string
}

type LockStateResponse struct {
	LockID      string
	Diagnostics []*tfprotov6.Diagnostic
}

type UnlockStateRequest struct {
	StateID string
	LockID  string
}

type UnlockStateResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type ReadStateBytesRequest struct {
	StateID string
}

type ReadStateBytesResponse struct {
	StateBytes  []byte
	Diagnostics []*tfprotov6.Diagnostic
}

type WriteStateBytesRequest struct {
	StateID    string
	StateBytes []byte
}

type WriteStateBytesResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}
