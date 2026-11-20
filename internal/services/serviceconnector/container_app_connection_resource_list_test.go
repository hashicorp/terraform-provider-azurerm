// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package serviceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccServiceConnectorContainerApp_list_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ContainerAppConnectionResource{}
	listResourceAddress := "azurerm_container_app_connection.list"
	containerAppId := containerapps.NewContainerAppID(data.Subscriptions.Primary, fmt.Sprintf("acctestRG-%d", data.RandomInteger), fmt.Sprintf("acctest-ca-%d", data.RandomInteger))
	connectionName := fmt.Sprintf("acctestserviceconnector0%d", data.RandomInteger)
	identity := map[string]knownvalue.Check{
		"subscription_id":     knownvalue.StringExact(containerAppId.SubscriptionId),
		"resource_group_name": knownvalue.StringExact(containerAppId.ResourceGroupName),
		"container_app_name":  knownvalue.StringExact(containerAppId.ContainerAppName),
		"name":                knownvalue.StringExact(connectionName),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.listQuery(containerAppId.ID(), false),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectIdentity(listResourceAddress, identity),
				},
			},
			{
				Query:  true,
				Config: r.listQuery(containerAppId.ID(), true),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectIdentity(listResourceAddress, identity),
					querycheck.ExpectResourceKnownValues(listResourceAddress, queryfilter.ByResourceIdentity(identity), []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("name"),
							KnownValue: knownvalue.StringExact(connectionName),
						},
						{
							Path:       tfjsonpath.New("container_app_id"),
							KnownValue: knownvalue.StringExact(containerAppId.ID()),
						},
						{
							Path:       tfjsonpath.New("scope"),
							KnownValue: knownvalue.StringExact(fmt.Sprintf("acctest-cont-%d", data.RandomInteger)),
						},
					}),
				},
			},
		},
	})
}

func (r ContainerAppConnectionResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_connection" "test" {
  count = 3

  name               = "acctestserviceconnector${count.index}%[2]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  scope              = azurerm_container_app.test.template[0].container[0].name

  authentication {
    type = "systemAssignedIdentity"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppConnectionResource) listQuery(containerAppId string, includeResource bool) string {
	return fmt.Sprintf(`
list "azurerm_container_app_connection" "list" {
  provider         = azurerm
  include_resource = %[2]t
  config {
    container_app_id = "%[1]s"
  }
}
`, containerAppId, includeResource)
}
