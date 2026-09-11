// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"context"
)

// NOTE: State store support is experimental and exposed without compatibility promises until
// these notices are removed.
type StateStore interface {
	// Metadata should return the full name of the state store, such
	// as examplecloud_store.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema should return the schema for this state store.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Initialize is called one time, prior to executing any state store RPCs (excluding offline validation) but after
	// the provider is configured, and is when Terraform sends the values the user specified in the state_store configuration
	// block to the provider. These are supplied in the InitializeRequest argument.
	//
	// As this method is always executed after provider configuration, data can be passed from [provider.ConfigureResponse.StateStoreData]
	// to [InitializeRequest.ProviderData]. This provider-level data along with the values from state store configuration are often used
	// to initialize an API client, which can be set to [InitializeResponse.StateStoreData], then eventually stored on the struct implementing
	// the StateStore interface in the [StateStoreWithConfigure.Configure] method.
	Initialize(context.Context, InitializeRequest, *InitializeResponse)

	// GetStates returns all state IDs for states persisted in the configured state store.
	GetStates(context.Context, GetStatesRequest, *GetStatesResponse)

	// DeleteState is called by Terraform to delete a state from the configured state store.
	DeleteState(context.Context, DeleteStateRequest, *DeleteStateResponse)

	// Lock is called by Terraform to acquire a lock prior to performing an operation that needs to write to state. If the [LockResponse.LockID] field
	// is a non-empty string, Terraform will call [StateStore.Unlock] once the operation has been completed.
	//
	// State stores that support locking are expected to handle concurrent clients by ensuring multiple locks cannot be acquired on the same state
	// simultaneously. The backing data store must be strongly consistent (i.e. a newly created lock is immediately visible to all clients) and some form
	// of concurrency control must be implemented when attempting to acquire a lock. An example of this would be creating a lock file with a conditional
	// write that would fail if the requested file already exists.
	Lock(context.Context, LockRequest, *LockResponse)

	// Unlock is called by Terraform to release a lock (previously acquired by [StateStore.Lock]) after an operation has been completed.
	//
	// This method is not called by Terraform if the state store returns an empty [LockResponse.LockID] from [StateStore.Lock].
	Unlock(context.Context, UnlockRequest, *UnlockResponse)

	// Read returns the given state as bytes from a state store.
	Read(context.Context, ReadRequest, *ReadResponse)

	// Write is called by Terraform to write state data to a given state ID in a state store.
	Write(context.Context, WriteRequest, *WriteResponse)
}

// StateStoreWithConfigure is an interface type that extends StateStore to
// include a method which the framework will automatically call so provider
// developers have the opportunity to setup any necessary provider-level data
// or clients in the StateStore type.
type StateStoreWithConfigure interface {
	StateStore

	// Configure enables provider-level data or clients to be set in the
	// provider-defined StateStore type.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
}

// StateStoreWithConfigValidators is an interface type that extends StateStore to include declarative validations.
//
// Declaring validation using this methodology simplifies implmentation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type StateStoreWithConfigValidators interface {
	StateStore

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// StateStoreWithValidateConfig is an interface type that extends StateStore to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single statestore. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type StateStoreWithValidateConfig interface {
	StateStore

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
