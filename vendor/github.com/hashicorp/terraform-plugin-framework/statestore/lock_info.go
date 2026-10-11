// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// LockInfo is a helper struct for representing lock data for a given state. The [NewLockInfo]
// function should be used to create this struct, which will populate [LockInfo.ID], [LockInfo.Who]
// and [LockInfo.Created] automatically. LockInfo structs are meant to be directly marshalled to JSON,
// stored, and unmarshalled from JSON for inspection.
//
// The [LockInfo.Equal] method can be used to compare locks.
//
// The [WorkspaceAlreadyLockedDiagnostic] function can create a formatted diagnostic that describes to
// Terraform that a new lock cannot be acquired as one already exists.
type LockInfo struct {
	// ID is a unique identifier for a lock. The [NewLockInfo] function should be
	// used to create a new LockInfo struct, where this ID will be generated.
	//
	// This field should be used to set [LockResponse.LockID].
	ID string `json:"id"`

	// Operation is the operation provided by the [LockRequest] (plan, apply, refresh, etc.)
	Operation string `json:"operation"`

	// Who is a string representing the <username>@<hostname> (if available).
	Who string `json:"who,omitempty"`

	// Created is the time that the lock was created, set by [NewLockInfo].
	Created time.Time `json:"created"`
}

// NewLockInfo will create a new [LockInfo] struct, which is a helper struct that
// can be marshalled to JSON and stored to represent an active lock for a state file.
func NewLockInfo(req LockRequest) LockInfo {
	return LockInfo{
		ID:        generateLockID(),
		Operation: req.Operation,
		Who:       generateLockInfoWho(),
		Created:   time.Now().UTC(),
	}
}

func (l LockInfo) Equal(other LockInfo) bool {
	if l.ID != other.ID {
		return false
	}
	if l.Operation != other.Operation {
		return false
	}
	if l.Who != other.Who {
		return false
	}
	if !l.Created.Equal(other.Created) {
		return false
	}
	return true
}

func (l LockInfo) String() string {
	return fmt.Sprintf(`Lock Info:
  ID:        %s
  Operation: %s
  Who:       %s
  Created:   %s
`, l.ID, l.Operation, l.Who, l.Created)
}

// WorkspaceAlreadyLockedDiagnostic returns an error diagnostic indicating that a lock is already held by another user
// for a given state ID (workspace).
func WorkspaceAlreadyLockedDiagnostic(req LockRequest, existingLock LockInfo) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Workspace Already Locked",
		fmt.Sprintf("%q workspace has a lock held by another client.\n\n%s", req.StateID, existingLock),
	)
}

var rngSource = rand.New(rand.NewSource(time.Now().UnixNano()))

// generateLockID returns a random UUID string
func generateLockID() string {
	buf := make([]byte, 16)
	rngSource.Read(buf)

	// The only error returned from this function is if the byte length isn't 16, so we can ignore it
	id, _ := uuid.FormatUUID(buf)

	return id
}

// generateLockInfoWho produces a string formatted as: <username>@<hostname>
//
// This is a best effort function and will return an empty string if any
// errors are encountered reading user/host information.
func generateLockInfoWho() string {
	userName := ""
	if userInfo, err := user.Current(); err == nil {
		userName = userInfo.Username
	}
	host, _ := os.Hostname()

	// This is just nice-to-have info, so if we experience an error here,
	// we'll just return any information we have.
	if len(host) == 0 {
		return userName
	}
	if len(userName) == 0 {
		return ""
	}

	return fmt.Sprintf("%s@%s", userName, host)
}
