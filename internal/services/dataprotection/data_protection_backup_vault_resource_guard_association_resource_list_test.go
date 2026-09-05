// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDataProtectionBackupVaultResourceGuardAssociation_listByBackupVaultID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_resource_guard_association", "test")
	r := DataProtectionBackupVaultResourceGuardAssociationResource{}
	listResourceAddress := "azurerm_data_protection_backup_vault_resource_guard_association.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringExact("DppResourceGuardProxy"),
							"backup_vault_name":   knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_data_protection_backup_vault_resource_guard_association" "list" {
  provider = azurerm
  config {
    data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
