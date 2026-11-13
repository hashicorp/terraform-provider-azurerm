// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_maintenance_configuration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccContainerAppEnvironmentMaintenanceConfiguration_listByContainerAppEnvironmentID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_maintenance_configuration", "test")
	r := ContainerAppEnvironmentMaintenanceConfigurationResource{}

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
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_container_app_environment_maintenance_configuration.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_container_app_environment_maintenance_configuration.list",
						map[string]knownvalue.Check{
							"name":                     knownvalue.StringExact("default"),
							"managed_environment_name": knownvalue.StringExact(fmt.Sprintf("acctest-CAEnv%d", data.RandomInteger)),
							"resource_group_name":      knownvalue.StringExact(fmt.Sprintf("acctestRG-CAE-%d", data.RandomInteger)),
							"subscription_id":          knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) basicQuery() string {
	return `
list "azurerm_container_app_environment_maintenance_configuration" "list" {
  provider = azurerm
  config {
    container_app_environment_id = azurerm_container_app_environment.test.id
  }
}
`
}
