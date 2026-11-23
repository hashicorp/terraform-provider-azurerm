// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
)

func flattenPrivilege(input *[]mongorbacs.Privilege) []Privilege {
	if input == nil {
		return []Privilege{}
	}

	result := make([]Privilege, 0, len(*input))
	for _, input := range *input {
		privilege := Privilege{
			Actions:  pointer.From(input.Actions),
			Resource: flattenResource(input.Resource),
		}

		result = append(result, privilege)
	}

	return result
}

func flattenResource(input *mongorbacs.PrivilegeResource) []ResourceModel {
	var result []ResourceModel
	if input == nil {
		return []ResourceModel{}
	}

	resource := ResourceModel{
		CollectionName: pointer.From(input.Collection),
		DbName:         pointer.From(input.Db),
	}

	return append(result, resource)
}

func flattenInheritedRoles(input *[]mongorbacs.Role) []string {
	if input == nil {
		return []string{}
	}

	result := make([]string, 0, len(*input))
	for _, v := range *input {
		result = append(result, pointer.From(v.Role))
	}

	return result
}
