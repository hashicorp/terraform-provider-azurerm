// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package quota_test

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

func TestAccQuotaGroup_list_basic(t *testing.T) {
	r := QuotaGroupResource{}
	listResourceAddress := "azurerm_quota_group.list"

	data := acceptance.BuildTestData(t, "azurerm_quota_group", "test")

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
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryIncludeResource(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r QuotaGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "test" {
  count = 3

  name                = "acctestqg${count.index}%d"
  management_group_id = data.azurerm_management_group.test.id
}
`, data.RandomInteger)
}

func (r QuotaGroupResource) basicQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_quota_group" "list" {
  provider = azurerm
  config {
    management_group_id = "/providers/Microsoft.Management/managementGroups/%s"
  }
}
`, data.Client().TenantID)
}

func (r QuotaGroupResource) basicQueryIncludeResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_quota_group" "list" {
  provider         = azurerm
  include_resource = true
  config {
    management_group_id = "/providers/Microsoft.Management/managementGroups/%s"
  }
}
`, data.Client().TenantID)
}
