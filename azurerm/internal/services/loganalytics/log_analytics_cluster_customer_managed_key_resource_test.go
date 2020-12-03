package loganalytics_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccLogAnalyticsClusterCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckLogAnalyticsClusterCustomerManagedKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogAnalyticsClusterCustomerManagedKey_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckLogAnalyticsClusterCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckLogAnalyticsClusterCustomerManagedKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Cluster Custoemr Managed Key not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsClusterID(strings.TrimRight(rs.Primary.ID, "/CMK"))
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: get on Log Analytics Cluster for CMK: %+v", err)
			}
		}
		if resp.ClusterProperties == nil || resp.ClusterProperties.KeyVaultProperties == nil {
			return fmt.Errorf("bad: Log Analytics Cluster has no Cutomer Managed Key Configured")
		}
		if resp.ClusterProperties.KeyVaultProperties.KeyVaultURI == nil || *resp.ClusterProperties.KeyVaultProperties.KeyVaultURI == "" {
			return fmt.Errorf("bad: Log Analytics Cluster Customer Managed Key is not configured")
		}
		if resp.ClusterProperties.KeyVaultProperties.KeyName == nil || *resp.ClusterProperties.KeyVaultProperties.KeyName == "" {
			return fmt.Errorf("bad: Log Analytics Cluster Customer Managed Key is not configured")
		}

		return nil
	}
}

func testCheckLogAnalyticsClusterCustomerManagedKeyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_cluster_customer_managed_key" {
			continue
		}
		id, err := parse.LogAnalyticsClusterID(strings.TrimRight(rs.Primary.ID, "/CMK"))
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: get on Log Analytics Cluster for CMK: %+v", err)
			}
		}
		if resp.ClusterProperties != nil && resp.ClusterProperties.KeyVaultProperties != nil {
			if resp.ClusterProperties.KeyVaultProperties.KeyName != nil || *resp.ClusterProperties.KeyVaultProperties.KeyName != "" {
				return fmt.Errorf("Bad: Log Analytics CLuster Customer Managed Key %q still present", *resp.ClusterProperties.KeyVaultProperties.KeyName)
			}
			if resp.ClusterProperties.KeyVaultProperties.KeyVaultURI != nil || *resp.ClusterProperties.KeyVaultProperties.KeyVaultURI != "" {
				return fmt.Errorf("Bad: Log Analytics CLuster Customer Managed Key Vault URI %q still present", *resp.ClusterProperties.KeyVaultProperties.KeyVaultURI)
			}
			if resp.ClusterProperties.KeyVaultProperties.KeyVersion != nil || *resp.ClusterProperties.KeyVaultProperties.KeyVersion != "" {
				return fmt.Errorf("Bad: Log Analytics CLuster Customer Managed Key Version %q still present", *resp.ClusterProperties.KeyVaultProperties.KeyVersion)
			}
		}
		return nil
	}
	return nil
}

func testAccLogAnalyticsClusterCustomerManagedKey_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_key_vault" "test" {
  name                = "vault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  soft_delete_enabled        = true
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}


resource "azurerm_key_vault_access_policy" "terraform" {
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

resource "azurerm_key_vault_key" "test" {
  name         = "key-%[3]s"
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

  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "unwrapkey",
    "wrapkey"
  ]

  tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
  object_id = azurerm_log_analytics_cluster.test.identity.0.principal_id

  depends_on = [azurerm_key_vault_access_policy.terraform]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccLogAnalyticsClusterCustomerManagedKey_complete(data acceptance.TestData) string {
	template := testAccLogAnalyticsClusterCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = azurerm_key_vault_key.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}

`, template)
}
