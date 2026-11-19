// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

type VaultGuardProxyModel struct {
	VaultId         string `tfschema:"vault_id"`
	ResourceGuardId string `tfschema:"resource_guard_id"`
}

func (r Resource) ModelObject() any {
	return &VaultGuardProxyModel{}
}
