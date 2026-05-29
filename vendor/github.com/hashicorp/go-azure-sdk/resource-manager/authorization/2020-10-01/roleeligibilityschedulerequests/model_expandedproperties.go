package roleeligibilityschedulerequests

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpandedProperties struct {
	Principal      *ExpandedPropertiesPrincipal      `json:"principal,omitempty"`
	RoleDefinition *ExpandedPropertiesRoleDefinition `json:"roleDefinition,omitempty"`
	Scope          *ExpandedPropertiesScope          `json:"scope,omitempty"`
}
