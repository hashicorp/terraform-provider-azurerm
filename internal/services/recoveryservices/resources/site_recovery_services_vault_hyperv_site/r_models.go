// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_services_vault_hyperv_site

type HyperVSiteModel struct {
	Name            string `tfschema:"name"`
	RecoveryVaultId string `tfschema:"recovery_vault_id"`
}

func (r Resource) ModelObject() any {
	return &HyperVSiteModel{}
}
