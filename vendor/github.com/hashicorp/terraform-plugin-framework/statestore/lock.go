// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import "github.com/hashicorp/terraform-plugin-framework/diag"

// LockRequest represents a request to lock a given state ([LockRequest.StateID]) in the state store.
// If the given state has already been locked, the provider should return an error diagnostic.
//
// The [LockInfo] struct can be used as a simple representation of a lock for storing and comparing
// persisted lock data by passing [LockRequest] to the [NewLockInfo] helper function.
type LockRequest struct {
	// StateID is the ID of the state to lock.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces .
	//
	// If the practitioner hasn't explicitly selected a workspace, StateID will be set to "default".
	StateID string

	// Operation represents the type of operation Terraform is running when requesting the lock (refresh, plan, apply, etc).
	Operation string
}

// LockResponse represents a response to an LockRequest. An instance of this response
// struct is supplied as an argument to the state store's Lock method, in which the provider
// should set values on the LockResponse as appropriate.
type LockResponse struct {
	// LockID is an opaque string representing a new lock that has been persisted in the configured state store
	// for a given state ([LockRequest.StateID]). LockID is determined by the provider and will be passed to
	// [UnlockRequest.LockID] to release the lock once the operation is complete.
	//
	// If the state store doesn't support locking or the current state store configuration is not setup for locking,
	// return the LockResponse with LockID unset and no diagnostics. This will inform Terraform the state store cannot be locked,
	// which will skip unlocking the state when the operation is complete.
	//
	// If the given state ([LockRequest.StateID]) has already been locked, the provider should return an error diagnostic with
	// information about the current lock. If the [LockInfo] struct is being used, this diagnostic can be created with the
	// [WorkspaceAlreadyLockedDiagnostic] function.
	LockID string

	// Diagnostics report errors or warnings related to locking a state in the given
	// state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
