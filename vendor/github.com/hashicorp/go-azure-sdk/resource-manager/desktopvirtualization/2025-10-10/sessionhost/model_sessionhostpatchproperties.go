package sessionhost

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionHostPatchProperties struct {
	AllowNewSession *bool   `json:"allowNewSession,omitempty"`
	AssignedUser    *string `json:"assignedUser,omitempty"`
	FriendlyName    *string `json:"friendlyName,omitempty"`
}
