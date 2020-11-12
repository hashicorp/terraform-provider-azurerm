package tests

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
	resource.ParallelTest(t, resource.TestCase{
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
			data.ImportStep("size_gb"), // not returned by the API
		},
	})
}

func TestAccAzureRMLogAnalyticsCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMLogAnalyticsCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
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
			data.ImportStep("size_gb"), // not returned by the API
			{
				Config: testAccAzureRMLogAnalyticsCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsClusterExists(data.ResourceName),
				),
			},
			data.ImportStep("size_gb"), // not returned by the API
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
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: log analytics Cluster %q does not exist", id.Name)
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
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
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

func testAccAzureRMLogAnalyticsCluster_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsCluster_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  soft_delete_enabled        = true
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault_access_policy.subscription]
}

resource "azurerm_key_vault_access_policy" "subscription" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "create",
    "delete",
    "get",
    "update",
    "list",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "unwrapkey",
    "wrapkey"
  ]

  tenant_id = azurerm_log_analytics_cluster.example.identity.0.tenant_id
  object_id = azurerm_log_analytics_cluster.example.identity.0.principal_id
}

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  key_vault_property {
    key_name      = azurerm_key_vault_key.test.name
    key_vault_uri = azurerm_key_vault.test.vault_uri
    key_version   = azurerm_key_vault_key.test.version
  }

  size_gb = 1100

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.RandomString, data.RandomInteger)
}
