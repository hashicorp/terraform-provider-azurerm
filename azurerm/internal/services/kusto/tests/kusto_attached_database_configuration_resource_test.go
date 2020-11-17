package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoAttachedDatabaseConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_attached_database_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoAttachedDatabaseConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoAttachedDatabaseConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoAttachedDatabaseConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKustoAttachedDatabaseConfiguration_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster1" {
  name                = "acctestkc1%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_cluster" "cluster2" {
  name                = "acctestkc2%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "followed_database" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster1.name
}

resource "azurerm_kusto_attached_database_configuration" "configuration1" {
  name                = "acctestka-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster1.name
  cluster_resource_id = azurerm_kusto_cluster.cluster2.id
  database_name       = "*"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMKustoAttachedDatabaseConfigurationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.AttachedDatabaseConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_attached_database_configuration" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMKustoAttachedDatabaseConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.AttachedDatabaseConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		configurationName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Attached Database Configuration: %s", configurationName)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Attached Database Configuration: %s", configurationName)
		}

		resp, err := client.Get(ctx, resourceGroup, clusterName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Attached Database Configuration %q (resource group: %q, cluster: %q) does not exist", configurationName, resourceGroup, clusterName)
			}

			return fmt.Errorf("Bad: Get on AttachedDatabaseConfigurationsClient: %+v", err)
		}

		return nil
	}
}
