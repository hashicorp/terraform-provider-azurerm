package rolemanagementpolicyassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyAssignmentProperties struct {
	Policy         *PolicyAssignmentPropertiesPolicy         `json:"policy,omitempty"`
	RoleDefinition *PolicyAssignmentPropertiesRoleDefinition `json:"roleDefinition,omitempty"`
	Scope          *PolicyAssignmentPropertiesScope          `json:"scope,omitempty"`
}
