// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_assignment

type KeyVaultManagedHSMRoleAssignmentModel struct {
	ManagedHSMID     string `tfschema:"managed_hsm_id"`
	Name             string `tfschema:"name"`
	Scope            string `tfschema:"scope"`
	RoleDefinitionId string `tfschema:"role_definition_id"`
	PrincipalId      string `tfschema:"principal_id"`
	ResourceId       string `tfschema:"resource_id"`
}

func (r Resource) ModelObject() any {
	return &KeyVaultManagedHSMRoleAssignmentModel{}
}
