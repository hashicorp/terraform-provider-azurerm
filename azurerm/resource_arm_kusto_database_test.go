package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoDatabase_basic(t *testing.T) {
	resourceName := "azurerm_kusto_database.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKustoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoDatabase_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabaseExists(resourceName),
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

func TestAccAzureRMKustoDatabase_softDeletePeriod(t *testing.T) {
	resourceName := "azurerm_kusto_database.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	preConfig := testAccAzureRMKustoDatabase_softDeletePeriod(ri, rs, testLocation())
	postConfig := testAccAzureRMKustoDatabase_softDeletePeriodUpdate(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKustoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "soft_delete_period", "P7D"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "soft_delete_period", "P31D"),
				),
			},
		},
	})
}

func TestAccAzureRMKustoDatabase_hotCachePeriod(t *testing.T) {
	resourceName := "azurerm_kusto_database.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	preConfig := testAccAzureRMKustoDatabase_hotCachePeriod(ri, rs, testLocation())
	postConfig := testAccAzureRMKustoDatabase_hotCachePeriodUpdate(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKustoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hot_cache_period", "P7D"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hot_cache_period", "P14DT12H"),
				),
			},
		},
	})
}

func testAccAzureRMKustoDatabase_basic(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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
`, rInt, location, rs, rInt)
}

func testAccAzureRMKustoDatabase_softDeletePeriod(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

  soft_delete_period = "P7D"
}
`, rInt, location, rs, rInt)
}

func testAccAzureRMKustoDatabase_softDeletePeriodUpdate(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

  soft_delete_period = "P31D"
}
`, rInt, location, rs, rInt)
}

func testAccAzureRMKustoDatabase_hotCachePeriod(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

  hot_cache_period = "P7D"
}
`, rInt, location, rs, rInt)
}

func testAccAzureRMKustoDatabase_hotCachePeriodUpdate(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

  hot_cache_period = "P14DT12H"
}
`, rInt, location, rs, rInt)
}

func testCheckAzureRMKustoDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).kusto.DatabasesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_database" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		name := rs.Primary.Attributes["name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMKustoDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		kustoDatabase := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Database: %s", kustoDatabase)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Database: %s", kustoDatabase)
		}

		client := testAccProvider.Meta().(*ArmClient).kusto.DatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, clusterName, kustoDatabase)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Database %q (resource group: %q, cluster: %q) does not exist", kustoDatabase, resourceGroup, clusterName)
			}

			return fmt.Errorf("Bad: Get on DatabasesClient: %+v", err)
		}

		return nil
	}
}
