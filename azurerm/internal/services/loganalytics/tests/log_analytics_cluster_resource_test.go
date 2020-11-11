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

func TestAccAzureRMoperationalinsightsCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMoperationalinsightsCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMoperationalinsightsCluster_requiresImport),
		},
	})
}

func TestAccAzureRMoperationalinsightsCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMoperationalinsightsCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMoperationalinsightsCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMoperationalinsightsCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMoperationalinsightsCluster_updateKeyVaultProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMoperationalinsightsCluster_updateKeyVaultProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMoperationalinsightsCluster_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMoperationalinsightsCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMoperationalinsightsCluster_updateSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMoperationalinsightsClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("operationalinsights Cluster not found: %s", resourceName)
		}
		id, err := parse.OperationalinsightsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: log_analytics Cluster %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on LogAnalytics.ClusterClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMoperationalinsightsClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_cluster" {
			continue
		}
		id, err := parse.OperationalinsightsClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on LogAnalytics.ClusterClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMoperationalinsightsCluster_template(data acceptance.TestData) string {
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

func testAccAzureRMoperationalinsightsCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMoperationalinsightsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMoperationalinsightsCluster_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMoperationalinsightsCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "import" {
  name = azurerm_log_analytics_cluster.test.name
  resource_group_name = azurerm_log_analytics_cluster.test.resource_group_name
  location = azurerm_log_analytics_cluster.test.location
}
`, config)
}

func testAccAzureRMoperationalinsightsCluster_complete(data acceptance.TestData) string {
	template := testAccAzureRMoperationalinsightsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  identity {
    type = ""
  }

  next_link = ""

  key_vault_property {
    key_name = ""
    key_vault_uri = ""
    key_version = ""
  }

  size_gb = 1100

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMoperationalinsightsCluster_updateKeyVaultProperties(data acceptance.TestData) string {
	template := testAccAzureRMoperationalinsightsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  identity {
    type = ""
  }

  next_link = ""

  key_vault_property {
    key_name = ""
    key_vault_uri = ""
    key_version = ""
  }

  size_gb = 1000

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMoperationalinsightsCluster_updateSku(data acceptance.TestData) string {
	template := testAccAzureRMoperationalinsightsCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  identity {
    type = ""
  }

  next_link = ""

  key_vault_property {
    key_name = ""
    key_vault_uri = ""
    key_version = ""
  }

  size_gb = 1000

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
