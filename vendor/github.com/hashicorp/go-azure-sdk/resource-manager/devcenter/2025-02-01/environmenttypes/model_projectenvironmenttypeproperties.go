package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectEnvironmentTypeProperties struct {
	CreatorRoleAssignment *ProjectEnvironmentTypeUpdatePropertiesCreatorRoleAssignment `json:"creatorRoleAssignment,omitempty"`
	DeploymentTargetId    *string                                                      `json:"deploymentTargetId,omitempty"`
	DisplayName           *string                                                      `json:"displayName,omitempty"`
	EnvironmentCount      *int64                                                       `json:"environmentCount,omitempty"`
	ProvisioningState     *ProvisioningState                                           `json:"provisioningState,omitempty"`
	Status                *EnvironmentTypeEnableStatus                                 `json:"status,omitempty"`
	UserRoleAssignments   *map[string]UserRoleAssignment                               `json:"userRoleAssignments,omitempty"`
}
