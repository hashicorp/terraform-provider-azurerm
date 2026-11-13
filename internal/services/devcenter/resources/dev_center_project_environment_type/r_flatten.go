// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_environment_type

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
)

func flattenDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}

func flattenDevCenterProjectEnvironmentTypeUserRoleAssignment(input *map[string]environmenttypes.UserRoleAssignment) []DevCenterProjectEnvironmentTypeUserRoleAssignment {
	results := make([]DevCenterProjectEnvironmentTypeUserRoleAssignment, 0)

	if input == nil {
		return results
	}

	for k, v := range *input {
		result := DevCenterProjectEnvironmentTypeUserRoleAssignment{
			UserId: k,
			Roles:  flattenDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(v.Roles),
		}

		results = append(results, result)
	}

	return results
}

func flattenDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}
