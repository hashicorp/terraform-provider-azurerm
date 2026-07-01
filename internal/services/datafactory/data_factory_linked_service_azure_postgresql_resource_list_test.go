// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

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

func TestAccLinkedServiceAzurePostgreSQL_list_basic(t *testing.T) {
	r := LinkedServiceAzurePostgreSQLResource{}
	listResourceAddress := "azurerm_data_factory_linked_service_azure_postgresql.list"

	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test1")

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
				Config: r.basicQueryByDataFactoryId(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r LinkedServiceAzurePostgreSQLResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  count = 3

  name                = "acctestadfpostgresql-${count.index}-%[1]d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "SystemAssignedManagedIdentity"
  server              = "acctest-dev-%[1]d.postgres.database.azure.com"
  port                = 5432
  database_name       = "testdatabase"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinkedServiceAzurePostgreSQLResource) basicQueryByDataFactoryId() string {
	return `
list "azurerm_data_factory_linked_service_azure_postgresql" "list" {
  provider = azurerm
  config {
    data_factory_id = "${azurerm_data_factory.test.id}"
  }
}
`
}
