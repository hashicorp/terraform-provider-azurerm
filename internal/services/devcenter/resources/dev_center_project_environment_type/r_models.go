// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_environment_type

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func (r Resource) ModelObject() any {
	return &DevCenterProjectEnvironmentTypeResourceModel{}
}

type DevCenterProjectEnvironmentTypeResourceModel struct {
	Name                       string                                              `tfschema:"name"`
	Location                   string                                              `tfschema:"location"`
	DevCenterProjectId         string                                              `tfschema:"dev_center_project_id"`
	DeploymentTargetId         string                                              `tfschema:"deployment_target_id"`
	CreatorRoleAssignmentRoles []string                                            `tfschema:"creator_role_assignment_roles"`
	Identity                   []identity.ModelSystemAssignedUserAssigned          `tfschema:"identity"`
	UserRoleAssignment         []DevCenterProjectEnvironmentTypeUserRoleAssignment `tfschema:"user_role_assignment"`
	Tags                       map[string]string                                   `tfschema:"tags"`
}
