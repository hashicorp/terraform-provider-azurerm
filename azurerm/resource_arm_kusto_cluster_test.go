package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoCluster_basic(t *testing.T) {
	resourceName := "azurerm_kusto_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoCluster_basic(ri, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(resourceName),
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

func TestAccAzureRMKustoCluster_withTags(t *testing.T) {
	resourceName := "azurerm_kusto_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	preConfig := testAccAzureRMKustoCluster_withTags(ri, rs, acceptance.Location())
	postConfig := testAccAzureRMKustoCluster_withTagsUpdate(ri, rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMKustoCluster_sku(t *testing.T) {
	resourceName := "azurerm_kusto_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	preConfig := testAccAzureRMKustoCluster_basic(ri, rs, acceptance.Location())
	postConfig := testAccAzureRMKustoCluster_skuUpdate(ri, rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Dev(No SLA)_Standard_D11_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_D11_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
				),
			},
		},
	})
}

func testAccAzureRMKustoCluster_basic(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, rInt, location, rs)
}

func testAccAzureRMKustoCluster_withTags(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test"
  }
}
`, rInt, location, rs)
}

func testAccAzureRMKustoCluster_withTagsUpdate(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, rInt, location, rs)
}

func testAccAzureRMKustoCluster_skuUpdate(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_D11_v2"
    capacity = 2
  }
}
`, rInt, location, rs)
}

func testCheckAzureRMKustoClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

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

func testCheckAzureRMKustoClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		kustoCluster := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto Cluster: %s", kustoCluster)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, kustoCluster)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Cluster %q (resource group: %q) does not exist", kustoCluster, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ClustersClient: %+v", err)
		}

		return nil
	}
}
