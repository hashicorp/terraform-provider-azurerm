package kusto_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAzureRMKustoClusterCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoClusterCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterWithCustomerManagedKeyExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_version"),
				),
			},
			data.ImportStep(),
			{
				// Delete the encryption settings resource and verify it is gone
				Config: testAccAzureRMKustoClusterCustomerManagedKey_template(data),
				Check: resource.ComposeTestCheckFunc(
					// Then ensure the encryption settings on the Kusto cluster
					// have been reverted to their default state
					testCheckAzureRMKustoClusterExistsWithoutCustomerManagedKey("azurerm_kusto_cluster.test"),
				),
			},
		},
	})
}

func TestAccAzureRMKustoClusterCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoClusterCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterWithCustomerManagedKeyExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_version"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKustoClusterCustomerManagedKey_requiresImport),
		},
	})
}

func TestAccAzureRMKustoClusterCustomerManagedKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoClusterCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterWithCustomerManagedKeyExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_version"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKustoClusterCustomerManagedKey_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterWithCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMKustoClusterWithCustomerManagedKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
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

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Cluster %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on kustoClustersClient: %+v", err)
		}

		if props := resp.ClusterProperties; props != nil {
			if encryption := props.KeyVaultProperties; encryption == nil {
				return fmt.Errorf("Kusto Cluster encryption properties not found: %s", resourceName)
			}
		}

		return nil
	}
}

func testCheckAzureRMKustoClusterExistsWithoutCustomerManagedKey(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
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

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto Cluster %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on kustoClustersClient: %+v", err)
		}

		if props := resp.ClusterProperties; props != nil {
			if encryption := props.KeyVaultProperties; encryption != nil {
				return fmt.Errorf("Kusto Cluster encryption properties still found: %s", resourceName)
			}
		}

		return nil
	}
}

func testAccAzureRMKustoClusterCustomerManagedKey_basic(data acceptance.TestData) string {
	template := testAccAzureRMKustoClusterCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id   = azurerm_kusto_cluster.test.id
  key_vault_id = azurerm_key_vault.test.id
  key_name     = azurerm_key_vault_key.first.name
  key_version  = azurerm_key_vault_key.first.version
}
`, template)
}

func testAccAzureRMKustoClusterCustomerManagedKey_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKustoClusterCustomerManagedKey_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "import" {
  cluster_id   = azurerm_kusto_cluster_customer_managed_key.test.cluster_id
  key_vault_id = azurerm_kusto_cluster_customer_managed_key.test.key_vault_id
  key_name     = azurerm_kusto_cluster_customer_managed_key.test.key_name
  key_version  = azurerm_kusto_cluster_customer_managed_key.test.key_version
}
`, template)
}

func testAccAzureRMKustoClusterCustomerManagedKey_updated(data acceptance.TestData) string {
	template := testAccAzureRMKustoClusterCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id   = azurerm_kusto_cluster.test.id
  key_vault_id = azurerm_key_vault.test.id
  key_name     = azurerm_key_vault_key.second.name
  key_version  = azurerm_key_vault_key.second.version
}
`, template)
}

func testAccAzureRMKustoClusterCustomerManagedKey_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  soft_delete_enabled      = true
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "cluster" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_kusto_cluster.test.identity.0.tenant_id
  object_id    = azurerm_kusto_cluster.test.identity.0.principal_id

  key_permissions = ["get", "unwrapkey", "wrapkey"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["get", "list", "create", "delete", "recover"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
