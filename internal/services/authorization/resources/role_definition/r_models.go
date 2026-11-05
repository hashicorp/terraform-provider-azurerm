// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_definition

type RoleDefinitionModel struct {
	RoleDefinitionId         string            `tfschema:"role_definition_id"`
	Name                     string            `tfschema:"name"`
	Scope                    string            `tfschema:"scope"`
	Description              string            `tfschema:"description"`
	Permissions              []PermissionModel `tfschema:"permissions"`
	AssignableScopes         []string          `tfschema:"assignable_scopes"`
	RoleDefinitionResourceId string            `tfschema:"role_definition_resource_id"`
}

type PermissionModel struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

func (r Resource) ModelObject() any {
	return &RoleDefinitionModel{}
}
