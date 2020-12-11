package kusto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccKustoClusterPrincipalAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_principal_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckKustoClusterPrincipalAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKustoClusterPrincipalAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckKustoClusterPrincipalAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckKustoClusterPrincipalAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_cluster_principal_assignment" {
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

func testCheckKustoClusterPrincipalAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Cluster Principal Assignment: %s", name)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no cluster found in state for Kusto Cluster Principal Assignment: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q) does not exist", name, resourceGroup, clusterName)
			}

			return fmt.Errorf("Bad: Get on ClusterPrincipalAssignmentsClient: %+v", err)
		}

		return nil
	}
}

func testAccKustoClusterPrincipalAssignment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_cluster_principal_assignment" "test" {
  name                = "acctestkdpa%d"
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.test.name

  tenant_id      = data.azurerm_client_config.current.tenant_id
  principal_id   = data.azurerm_client_config.current.client_id
  principal_type = "App"
  role           = "AllDatabasesViewer"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
