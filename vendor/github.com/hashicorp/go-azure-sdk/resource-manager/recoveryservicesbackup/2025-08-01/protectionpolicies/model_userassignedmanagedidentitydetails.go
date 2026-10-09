package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedManagedIdentityDetails struct {
	IdentityArmId                  *string                         `json:"identityArmId,omitempty"`
	IdentityName                   *string                         `json:"identityName,omitempty"`
	UserAssignedIdentityProperties *UserAssignedIdentityProperties `json:"userAssignedIdentityProperties,omitempty"`
}
