// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

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

func TestAccContainerAppSessionPool_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "testlist")
	r := ContainerAppSessionPoolResource{}

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
				Config: r.subscriptionListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_container_app_session_pool.list", 2),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_container_app_session_pool.list", 2),
				},
			},
		},
	})
}

func (r ContainerAppSessionPoolResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_session_pool" "test" {
  count = 2

  name                    = "acctestcasp${count.index}%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  container_type          = "PythonLTS"
  max_concurrent_sessions = 5

  pool_management_type       = "Dynamic"
  lifecycle_type             = "Timed"
  cooldown_period_in_seconds = 300
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppSessionPoolResource) subscriptionListQuery() string {
	return `
list "azurerm_container_app_session_pool" "list" {
  provider = azurerm
}
`
}

func (r ContainerAppSessionPoolResource) resourceGroupListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_container_app_session_pool" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-CASP-%[1]d"
  }
}
`, data.RandomInteger)
}
