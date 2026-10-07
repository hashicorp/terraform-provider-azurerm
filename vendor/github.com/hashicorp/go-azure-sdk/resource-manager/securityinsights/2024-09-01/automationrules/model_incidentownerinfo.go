package automationrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IncidentOwnerInfo struct {
	AssignedTo        *string    `json:"assignedTo,omitempty"`
	Email             *string    `json:"email,omitempty"`
	ObjectId          *string    `json:"objectId,omitempty"`
	OwnerType         *OwnerType `json:"ownerType,omitempty"`
	UserPrincipalName *string    `json:"userPrincipalName,omitempty"`
}
