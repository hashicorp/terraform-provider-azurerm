// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

type RoleAssignmentModel struct {
	Name                               string `tfschema:"name"`
	PrincipalId                        string `tfschema:"principal_id"`
	RoleDefinitionId                   string `tfschema:"role_definition_id"`
	RoleDefinitionName                 string `tfschema:"role_definition_name"`
	SkipServicePrincipalAadCheck       bool   `tfschema:"skip_service_principal_aad_check"`
	DelegatedManagedIdentityResourceId string `tfschema:"delegated_managed_identity_resource_id"`
	Description                        string `tfschema:"description"`
	Condition                          string `tfschema:"condition"`
	ConditionVersion                   string `tfschema:"condition_version"`
	PrincipalType                      string `tfschema:"principal_type"`
}

func (r Resource) ModelObject() any {
	return &RoleAssignmentModel{}
}
