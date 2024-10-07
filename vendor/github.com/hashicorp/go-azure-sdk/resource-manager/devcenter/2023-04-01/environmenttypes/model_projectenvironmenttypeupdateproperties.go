package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectEnvironmentTypeUpdateProperties struct {
	CreatorRoleAssignment *ProjectEnvironmentTypeUpdatePropertiesCreatorRoleAssignment `json:"creatorRoleAssignment,omitempty"`
	DeploymentTargetId    *string                                                      `json:"deploymentTargetId,omitempty"`
	Status                *EnvironmentTypeEnableStatus                                 `json:"status,omitempty"`
	UserRoleAssignments   *map[string]UserRoleAssignment                               `json:"userRoleAssignments,omitempty"`
}
