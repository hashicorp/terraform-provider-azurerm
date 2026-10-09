// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	dataplane "github.com/jackofallops/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

type Role struct {
	RoleId         string `tfschema:"role_id"`
	OrganizationId string `tfschema:"organization_id"`
}

func ConvertToRoleAssignments(input []Role) *[]dataplane.RoleAssignment {
	if input == nil {
		return nil
	}

	results := make([]dataplane.RoleAssignment, 0)
	for _, item := range input {
		results = append(results, dataplane.RoleAssignment{
			Organization: pointer.To(item.OrganizationId),
			Role:         pointer.To(item.RoleId),
		})
	}
	return &results
}

func ConvertFromRoleAssignments(input *[]dataplane.RoleAssignment) []Role {
	if input == nil {
		return nil
	}

	results := make([]Role, 0)
	for _, item := range *input {
		obj := Role{}
		if item.Organization != nil {
			obj.OrganizationId = *item.Organization
		}
		if item.Role != nil {
			obj.RoleId = *item.Role
		}
		results = append(results, obj)
	}
	return results
}
