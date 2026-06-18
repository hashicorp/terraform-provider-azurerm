package synapseroledefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseRoleDefinition struct {
	AvailabilityStatus *string                  `json:"availabilityStatus,omitempty"`
	Description        *string                  `json:"description,omitempty"`
	Id                 *string                  `json:"id,omitempty"`
	IsBuiltIn          *bool                    `json:"isBuiltIn,omitempty"`
	Name               *string                  `json:"name,omitempty"`
	Permissions        *[]SynapseRbacPermission `json:"permissions,omitempty"`
	Scopes             *[]string                `json:"scopes,omitempty"`
}
