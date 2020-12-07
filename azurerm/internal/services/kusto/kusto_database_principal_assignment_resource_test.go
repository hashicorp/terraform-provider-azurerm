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

func TestAccAzureRMKustoDatabasePrincipalAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoDatabasePrincipalAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoDatabasePrincipalAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabasePrincipalAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKustoDatabasePrincipalAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoDatabasePrincipalAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoDatabasePrincipalAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabasePrincipalAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKustoDatabasePrincipalAssignment_requiresImport),
		},
	})
}

func testCheckAzureRMKustoDatabasePrincipalAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_database_principal_assignment" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
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

func testCheckAzureRMKustoDatabasePrincipalAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Database Principal Assignment: %s", name)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no cluster found in state for Kusto Database Principal Assignment: %s", name)
		}

		databaseName, hasDatabaseName := rs.Primary.Attributes["database_name"]
		if !hasDatabaseName {
			return fmt.Errorf("Bad: no database found in state for Kusto Database Principal Assignment: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q) does not exist", name, resourceGroup, clusterName, databaseName)
			}

			return fmt.Errorf("Bad: Get on DatabasePrincipalAssignmentsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKustoDatabasePrincipalAssignment_basic(data acceptance.TestData) string {
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

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_database_principal_assignment" "test" {
  name                = "acctestkdpa%d"
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  tenant_id      = data.azurerm_client_config.current.tenant_id
  principal_id   = data.azurerm_client_config.current.client_id
  principal_type = "App"
  role           = "Viewer"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKustoDatabasePrincipalAssignment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKustoDatabasePrincipalAssignment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_database_principal_assignment" "import" {
  name                = azurerm_kusto_database_principal_assignment.test.name
  resource_group_name = azurerm_kusto_database_principal_assignment.test.resource_group_name
  cluster_name        = azurerm_kusto_database_principal_assignment.test.cluster_name
  database_name       = azurerm_kusto_database_principal_assignment.test.database_name

  tenant_id      = azurerm_kusto_database_principal_assignment.test.tenant_id
  principal_id   = azurerm_kusto_database_principal_assignment.test.principal_id
  principal_type = azurerm_kusto_database_principal_assignment.test.principal_type
  role           = azurerm_kusto_database_principal_assignment.test.role
}
`, template)
}
