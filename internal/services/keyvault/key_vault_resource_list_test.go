// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccKeyVault_list_basic(t *testing.T) {
	r := KeyVaultResource{}
	listResourceAddress := "azurerm_key_vault.list"

	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test0")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r KeyVaultResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  count = 3

  name                       = "acctestkv${count.index}%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r KeyVaultResource) basicQuery() string {
	return `
list "azurerm_key_vault" "list" {
  provider = azurerm
  config {}
}
`
}

func (r KeyVaultResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_key_vault" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}
`, data.RandomInteger)
}
