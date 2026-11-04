// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import "github.com/hashicorp/terraform-plugin-framework/diag"

// UnlockRequest represents a request to release a lock ([UnlockRequest.LockID]) for a given state ([UnlockRequest.StateID])
// in the state store.
//
// The provider should return a diagnostic if a lock doesn't exist in the state store for the given state, or if the lock ID
// for the given state doesn't match [UnlockRequest.LockID]. These scenarios likely indicate either a bug in the [StateStore.Lock] method
// or the underlying storage mechanism where multiple concurrent clients were able to acquire a lock on the same state.
type UnlockRequest struct {
	// StateID is the ID of the state to unlock.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces .
	//
	// If the practitioner hasn't explicitly selected a workspace, StateID will be set to "default".
	StateID string

	// LockID is the ID of the lock to be released (unlocked) for a given state in the configured state store.
	// This is the same value that is returned when originally acquiring the lock from the [StateStore.Lock] method,
	// i.e. the [LockResponse.LockID] field.
	//
	// The provider should return a diagnostic if a lock doesn't exist in the state store for the given state, or if the lock ID
	// for the given state doesn't match LockID. These scenarios likely indicate either a bug in the [StateStore.Lock] method
	// or the underlying storage mechanism where multiple concurrent clients were able to acquire a lock on the same state.
	LockID string
}

// UnlockResponse represents a response to an UnlockRequest. An instance of this response
// struct is supplied as an argument to the state store's Unlock method, in which the provider
// should set values on the UnlockResponse as appropriate.
type UnlockResponse struct {
	// Diagnostics report errors or warnings related to unlocking a state in the given
	// state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
