package onlineendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartialManagedServiceIdentity struct {
	Type                   *ManagedServiceIdentityType `json:"type,omitempty"`
	UserAssignedIdentities *map[string]interface{}     `json:"userAssignedIdentities,omitempty"`
}
