package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoDatabasePrincipal_basic(t *testing.T) {
	resourceName := "azurerm_kusto_database_principal.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKustoDatabasePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoDatabasePrincipal_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabasePrincipalExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMKustoDatabasePrincipalDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Kusto.DatabasesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_database_principal" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		role := rs.Primary.Attributes["role"]
		fqn := rs.Primary.Attributes["fully_qualified_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		found := false
		if principals := resp.Value; principals != nil {
			for _, currPrincipal := range *principals {
				// kusto database principals are unique when looked at with fqn and role
				if string(currPrincipal.Role) == role && currPrincipal.Fqn != nil && *currPrincipal.Fqn == fqn {
					found = true
					break
				}
			}
		}
		if found {
			return fmt.Errorf("Kusto Database Principal %q still exists", fqn)
		}

		return nil
	}

	return nil
}

func testCheckAzureRMKustoDatabasePrincipalExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		role := rs.Primary.Attributes["role"]
		fqn := rs.Primary.Attributes["fully_qualified_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Database Principal: %s", fqn)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no cluster name found in state for Kusto Database Principal: %s", fqn)
		}

		databaseName, hasDatabaseName := rs.Primary.Attributes["database_name"]
		if !hasDatabaseName {
			return fmt.Errorf("Bad: no datbase name found in state for Kusto Database Principal: %s", fqn)
		}

		client := testAccProvider.Meta().(*ArmClient).Kusto.DatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Database %q (resource group: %q, cluster: %q) does not exist", fqn, resourceGroup, clusterName)
			}

			return fmt.Errorf("Bad: Get on DatabasesClient: %+v", err)
		}

		found := false
		if principals := resp.Value; principals != nil {
			for _, currPrincipal := range *principals {
				// kusto database principals are unique when looked at with fqn and role
				if string(currPrincipal.Role) == role && currPrincipal.Fqn != nil && *currPrincipal.Fqn == fqn {
					found = true
					break
				}
			}
		}
		if !found {
			return fmt.Errorf("Unable to find Kusto Database Principal %q", fqn)
		}

		return nil
	}
}

func testAccAzureRMKustoDatabasePrincipal_basic(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}


resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
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
  cluster_name        = azurerm_kusto_cluster.cluster.name
}

resource "azurerm_kusto_database_principal" "test" {
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.cluster.name
  database_name       = azurerm_kusto_database.test.name

  role      = "Viewer"
  type      = "App"
  client_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.client_id
}
`, rInt, location, rs, rInt)
}
