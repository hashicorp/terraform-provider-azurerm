// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

import (
	"context"
	"iter"
)

// StateStoreMetadata describes metadata for a state store in the GetMetadata RPC.
type StateStoreMetadata struct {
	// TypeName is the name of the state store.
	TypeName string
}

// StateStoreServer is an interface containing the methods a state store implementation needs to fill.
//
// NOTE: State store support is experimental and exposed without compatibility promises until
// these notices are removed.
type StateStoreServer interface {
	// ValidateStateStoreConfig performs configuration validation for the state store.
	ValidateStateStoreConfig(context.Context, *ValidateStateStoreConfigRequest) (*ValidateStateStoreConfigResponse, error)

	// ConfigureStateStore is called to pass the user-specified state store configuration to the provider, typically to store
	// and reference in future RPC calls.
	ConfigureStateStore(context.Context, *ConfigureStateStoreRequest) (*ConfigureStateStoreResponse, error)

	// ReadStateBytes returns a stream of byte chunks, for a given state, from a state store. The size of the byte chunks
	// are negotiated between Terraform and the provider in the ConfigureStateStore RPC.
	ReadStateBytes(context.Context, *ReadStateBytesRequest) (*ReadStateBytesStream, error)

	// WriteStateBytes receives a stream of byte chunks, for a given state, from Terraform to persist. The size of the
	// byte chunks are negotiated between Terraform and the provider in the ConfigureStateStore RPC.
	WriteStateBytes(context.Context, *WriteStateBytesStream) (*WriteStateBytesResponse, error)

	// GetStates returns a list of all states (i.e. CE workspaces) managed by a given state store.
	GetStates(context.Context, *GetStatesRequest) (*GetStatesResponse, error)

	// DeleteState instructs a given state store to delete a specific state. (i.e. a CE workspace)
	DeleteState(context.Context, *DeleteStateRequest) (*DeleteStateResponse, error)

	// LockState instructs a given state store to lock a specific state. (i.e. a CE workspace)
	LockState(context.Context, *LockStateRequest) (*LockStateResponse, error)

	// UnlockState instructs a given state store to unlock a specific state. (i.e. a CE workspace)
	UnlockState(context.Context, *UnlockStateRequest) (*UnlockStateResponse, error)
}

// ValidateStateStoreConfigRequest is the request Terraform sends when it
// wants to validate a state store's configuration.
type ValidateStateStoreConfigRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// Config is the configuration the user supplied for the state store. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	Config *DynamicValue
}

// ValidateStateStoreConfigResponse is the response from the provider about
// the validity of a state store's configuration.
type ValidateStateStoreConfigResponse struct {
	// Diagnostics report errors or warnings related to the given
	// configuration. Returning an empty slice indicates a successful
	// validation with no warnings or errors generated.
	Diagnostics []*Diagnostic
}

// ConfigureStateStoreRequest is the request Terraform sends to supply the
// provider with information about what the user entered in the state store's
// configuration block as well as negotiate capabilities with the state store implementation.
// This allows the provider to store and reference this information in future RPC calls.
type ConfigureStateStoreRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// Config is the configuration the user supplied for the state store. This
	// information should usually be persisted to the underlying type
	// that's implementing the StateStoreServer interface, for use in later
	// RPC requests. See the documentation on `DynamicValue` for more
	// information about safely accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	Config *DynamicValue

	// Capabilities defines protocol features for the ConfigureStateStore RPC,
	// such as forward-compatible Terraform behavior changes and size limitations
	// for the client.
	Capabilities *ConfigureStateStoreClientCapabilities
}

// ConfigureStateStoreResponse is the response from the provider about the
// state store configuration that Terraform supplied for the provider and
// the state store's capabilities.
type ConfigureStateStoreResponse struct {
	// Diagnostics report errors or warnings related to the state store's
	// configuration. Returning an empty slice indicates success, with no
	// errors or warnings generated.
	Diagnostics []*Diagnostic

	// Capabilities defines protocol features for the ConfigureStateStore RPC,
	// such as forward-compatible Terraform behavior changes and size limitations
	// for the server (state store implementation).
	Capabilities *StateStoreServerCapabilities
}

// StateStoreServerCapabilities allows providers to communicate supported
// protocol features related to state stores, such as forward-compatible
// Terraform behavior changes. StateStoreServerCapabilities is also used to communicate
// the provider chosen chunk size for future state store operations with Terraform.
type StateStoreServerCapabilities struct {
	// ChunkSize is the provider-chosen size of state byte chunks that will be sent between Terraform and
	// the provider in the ReadStateBytes and WriteStateBytes RPC calls. In most cases, provider implementations
	// can simply accept the provided ChunkSize from the client [ConfigureStateStoreRequest.Capabilities] request,
	// which for Terraform, will default to 8 MB.
	ChunkSize int64
}

// ReadStateBytesRequest is the request Terraform sends when Terraform needs to read state
// from a configured state store.
type ReadStateBytesRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// StateID is the ID of the state to read.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateID string
}

// ReadStateBytesStream represents a streaming response to a ReadStateBytesRequest. The provider should
// set a Chunks iterator function that pushes state byte chunks of type ReadStateByteChunk.
type ReadStateBytesStream struct {
	// Chunks is an iterator for returning state bytes, which can be chunked into smaller pieces if the overall size of the state
	// exceeds the max chunk size negotiated during ConfigureStateStore. Each chunk will be immediately sent to Terraform when
	// pushed to the iterator.
	//
	// Diagnostics that occur during the read operation can be returned by sending a chunk
	// with [ReadStateByteChunk.Diagnostics] set. Chunk sizes should not exceed the returned
	// ChunkSize capability during the ConfigureStateStore RPC.
	Chunks iter.Seq[ReadStateByteChunk]
}

// ReadStateByteChunk is a chunk of state byte data streamed from the provider to Terraform during the ReadStateBytes RPC.
//
// Diagnostics that occur during the read operation can be returned by setting the
// [ReadStateByteChunk.Diagnostics] field with no StateByteChunk data. The chunk size
// should not exceed the returned ChunkSize capability during the ConfigureStateStore RPC.
type ReadStateByteChunk struct {
	// StateByteChunk contains all of the necessary information about a chunk of state byte data.
	StateByteChunk

	// Diagnostics report errors or warnings related to retrieving the
	// state represented by the given state ID. Returning an empty slice
	// indicates success, with no errors or warnings generated.
	//
	// As diagnostics can be returned with each chunk, the error or warning can
	// be related to either the entire state or state store (returned with the first chunk), or
	// related to a specific chunk of data in the state.
	Diagnostics []*Diagnostic
}

// StateByteChunk represents a chunk of state byte data.
type StateByteChunk struct {
	// Bytes represents the state data for this chunk.
	Bytes []byte

	// TotalLength is the overall size of all of the state byte chunks that will be sent or received
	// for a given state.
	TotalLength int64

	// Range represents the start and end location of the [Bytes] for this chunk, in the context of
	// all of the chunks sent or received for a given state.
	Range StateByteRange
}

// StateByteRange represents the start and end location for a [StateByteChunk] in the context of
// all of the chunks sent or received for a given state.
type StateByteRange struct {
	// Start is the starting byte index for a [StateByteChunk].
	Start int64

	// End is the ending byte index for a [StateByteChunk].
	End int64
}

// WriteStateBytesStream represents a stream of state byte data sent from Terraform to the provider when it requests state be
// written to a configured state store. The iterator will continue producing [WriteStateBytesChunk] until all chunks have been
// sent from Terraform.
type WriteStateBytesStream struct {
	// Chunks represents a stream of state byte data chunks produced by Terraform. Consumers of this iterator should
	// always check the diagnostics before reading the [WriteStateBytesChunk] data. Diagnostics produced by this iterator would
	// indicate either invalid chunk data or GRPC communication errors.
	Chunks iter.Seq2[*WriteStateBytesChunk, []*Diagnostic]
}

// WriteStateBytesChunk is a chunk of state byte data, received from Terraform to be persisted. The chunk size should not exceed the returned
// ChunkSize capability during the ConfigureStateStore RPC.
type WriteStateBytesChunk struct {
	// Meta contains additional information about the WriteStateBytes RPC call necessary for routing the request in a provider.
	//
	// This field is only set with the first chunk sent from Terraform.
	Meta *WriteStateChunkMeta

	// StateByteChunk contains all of the necessary information about a chunk of state byte data.
	StateByteChunk
}

// WriteStateChunkMeta represents the additional information communicated during a WriteStateBytes RPC call. This information
// is required by providers to route the request and persist the state byte data with the corresponding state ID.
type WriteStateChunkMeta struct {
	// TypeName is the name of the state store.
	TypeName string

	// StateID is the ID of the state to write.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateID string
}

// WriteStateBytesResponse is the response from the provider after writing all chunks streamed from
// Terraform in a [WriteStateBytesStream] to a given state.
type WriteStateBytesResponse struct {
	// Diagnostics report errors or warnings related to writing the
	// state represented by the given state ID. Returning an empty slice
	// indicates success, with no errors or warnings generated.
	Diagnostics []*Diagnostic
}

// GetStatesRequest is the request Terraform sends when it needs to retrieve all peristed states
// in a configured state store.
type GetStatesRequest struct {
	// TypeName is the name of the state store.
	TypeName string
}

// GetStatesResponse is the response from the provider when retrieving all peristed states in a configured
// state store.
type GetStatesResponse struct {
	// StateIDs is a list of the states persisted to the configured state store.
	//
	// Typically these are the names of the Terraform workspaces the practitioner has persisted
	// in the state store: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateIDs []string

	// Diagnostics report errors or warnings related to retrieving all
	// state IDs stored by the configured state store. Returning an empty
	// slice indicates success, with no errors or warnings generated.
	Diagnostics []*Diagnostic
}

// DeleteStateRequest is the request Terraform sends when it needs to delete a given state
// in a configured state store.
type DeleteStateRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// StateID is the ID of the state to delete.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateID string
}

// DeleteStateResponse is the response from the provider after deleting a state in the configured state store.
type DeleteStateResponse struct {
	// Diagnostics report errors or warnings related to deleting the state
	// represented by the given state ID. Returning an empty slice indicates
	// success, with no errors or warnings generated.
	Diagnostics []*Diagnostic
}

// LockStateRequest is the request Terraform sends when it needs to lock a given state in a configured state store.
type LockStateRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// StateID is the ID of the state to lock.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateID string

	// Operation is a string representing the type of operation Terraform is currently running. (refresh, plan, apply, etc.)
	Operation string
}

// LockStateResponse is the response from the provider after locking a state in a configured state store.
type LockStateResponse struct {
	// LockID is an opaque string representing a new lock that has been persisted in the configured state store
	// for a given state. LockID is determined by the provider and will be passed to [UnlockStateRequest.LockID]
	// to release the lock.
	//
	// If the lock already exists, the provider should return an error diagnostic.
	LockID string

	// Diagnostics report errors or warnings related to locking the state
	// represented by the given state ID. Returning an empty slice indicates
	// success, with no errors or warnings generated.
	Diagnostics []*Diagnostic
}

// UnlockStateRequest is the request Terraform sends when it needs to unlock a given state in a configured state store.
type UnlockStateRequest struct {
	// TypeName is the name of the state store.
	TypeName string

	// StateID is the ID of the state to unlock.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces
	StateID string

	// LockID is the ID of the lock to be released (unlocked) for a given state in the configured state store.
	// This is the same value that is returned when originally acquiring the lock from the LockState RPC call,
	// i.e. the [LockStateResponse.LockID] field.
	LockID string
}

// UnlockStateResponse is the response from the provider after unlocking a state in the configured state store.
type UnlockStateResponse struct {
	// Diagnostics report errors or warnings related to unlocking the state
	// represented by the given state ID. Returning an empty slice indicates
	// success, with no errors or warnings generated.
	Diagnostics []*Diagnostic
}
