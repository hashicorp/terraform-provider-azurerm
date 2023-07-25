package managementlocks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementLockProperties struct {
	Level  LockLevel              `json:"level"`
	Notes  *string                `json:"notes,omitempty"`
	Owners *[]ManagementLockOwner `json:"owners,omitempty"`
}
