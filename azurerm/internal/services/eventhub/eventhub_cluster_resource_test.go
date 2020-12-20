package eventhub_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
)

func TestAccEventHubCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccEventHubCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccEventHubCluster_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccEventHubCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckEventHubClusterDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_cluster" {
			continue
		}

		id, err := parse.ClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("EventHub Cluster still exists:\n%#v", resp.ClusterProperties)
		}
	}

	return nil
}

func testCheckEventHubClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on clustersClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Event Hub Cluster %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccEventHubCluster_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctesteventhubclusTER-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccEventHubCluster_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctesteventhubclusTER-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated_1"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
