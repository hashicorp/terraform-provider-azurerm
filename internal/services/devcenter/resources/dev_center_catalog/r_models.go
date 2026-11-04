// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

type DevCenterCatalogsResourceModel struct {
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	DevCenterID       string                   `tfschema:"dev_center_id"`
	CatalogGitHub     []CatalogPropertiesModel `tfschema:"catalog_github"`
	CatalogAdoGit     []CatalogPropertiesModel `tfschema:"catalog_adogit"`
}

type CatalogPropertiesModel struct {
	URI            string `tfschema:"uri"`
	Branch         string `tfschema:"branch"`
	KeyVaultKeyUrl string `tfschema:"key_vault_key_url"`
	Path           string `tfschema:"path"`
}

func (r Resource) ModelObject() any {
	return &DevCenterCatalogsResourceModel{}
}
