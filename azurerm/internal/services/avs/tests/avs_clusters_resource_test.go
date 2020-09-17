package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

func TestAccAzureRMavsCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMavsCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsClusterExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMavsCluster_requiresImport),
		},
	})
}

func TestAccAzureRMavsCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMavsCluster_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMavsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMavsClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("avs Cluster not found: %s", resourceName)
		}
		id, err := parse.AvsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Avs Cluster %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Avs.ClusterClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMavsClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_avs_cluster" {
			continue
		}
		id, err := parse.AvsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Avs.ClusterClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMavsCluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-avs-%[1]d"
  location = "%[2]s"
}

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-apc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    cluster_size = 3
  }
  network_block = "192.168.48.0/22"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMavsCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMavsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_cluster" "test" {
  name             = "acctest-ac-%d"
  private_cloud_id = azurerm_avs_private_cloud.test.id
  cluster_size     = 3
  sku_name         = "av36t"
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsCluster_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMavsCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_cluster" "import" {
  name             = azurerm_avs_cluster.test.name
  private_cloud_id = azurerm_avs_cluster.test.private_cloud_id
  cluster_size     = azurerm_avs_cluster.test.cluster_size
  sku_name         = azurerm_avs_cluster.test.sku_name
}
`, config)
}

func testAccAzureRMavsCluster_update(data acceptance.TestData) string {
	template := testAccAzureRMavsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_cluster" "test" {
  name               = "acctest-ac-%d"
  private_cloud_name = azurerm_avs_private_cloud.test.id
  cluster_size       = 4
  sku_name           = "av36t"
}
`, template, data.RandomInteger)
}
