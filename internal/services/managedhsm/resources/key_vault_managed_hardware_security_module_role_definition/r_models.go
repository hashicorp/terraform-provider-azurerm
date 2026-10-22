// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_definition

type KeyVaultMHSMRoleDefinitionModel struct {
	ManagedHSMID      string       `tfschema:"managed_hsm_id"`
	Name              string       `tfschema:"name"`
	RoleName          string       `tfschema:"role_name"`
	Description       string       `tfschema:"description"`
	Permission        []Permission `tfschema:"permission"`
	RoleType          string       `tfschema:"role_type"`
	ResourceManagerId string       `tfschema:"resource_manager_id"`
}

func (r Resource) ModelObject() any {
	return &KeyVaultMHSMRoleDefinitionModel{}
}
