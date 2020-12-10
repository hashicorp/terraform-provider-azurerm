package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLogAnalyticsCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsCluster_resize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsCluster_resize(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsClusterExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsCluster_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("log analytics Cluster not found: %s", resourceName)
		}
		id, err := parse.LogAnalyticsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: log analytics Cluster %q does not exist", id.ClusterName)
			}
			return fmt.Errorf("bad: Get on LogAnalytics.ClusterClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMLogAnalyticsClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_cluster" {
			continue
		}
		id, err := parse.LogAnalyticsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on LogAnalytics.ClusterClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMLogAnalyticsCluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMLogAnalyticsCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsCluster_resize(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size_gb             = 1100

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsCluster_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMLogAnalyticsCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "import" {
  name                = azurerm_log_analytics_cluster.test.name
  resource_group_name = azurerm_log_analytics_cluster.test.resource_group_name
  location            = azurerm_log_analytics_cluster.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, config)
}
