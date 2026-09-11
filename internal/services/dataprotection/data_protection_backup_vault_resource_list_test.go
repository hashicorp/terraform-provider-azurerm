package dataprotection_test

import (
	"context"
	"fmt"
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

func TestAccDataProtectionBackupVault_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "testlist1")
	r := DataProtectionBackupVaultResource{}
	listResourceAddress := "azurerm_data_protection_backup_vault.list"

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
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r DataProtectionBackupVaultResource) basicQuery() string {
	return `
list "azurerm_data_protection_backup_vault" "list" {
  provider = azurerm
  config {}
}
`
}

func (r DataProtectionBackupVaultResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_data_protection_backup_vault" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}

func (r DataProtectionBackupVaultResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_vault" "test" {
  count               = 3
  name                = "acctest-bv-${count.index}-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
`, r.template(data), data.RandomInteger)
}
