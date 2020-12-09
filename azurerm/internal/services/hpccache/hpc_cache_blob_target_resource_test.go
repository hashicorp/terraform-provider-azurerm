package hpccache_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccHPCCacheBlobTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccHPCCacheBlobTarget_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccHPCCacheBlobTarget_namespace(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccHPCCacheBlobTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccHPCCacheBlobTarget_requiresImport),
		},
	})
}

func testCheckHPCCacheBlobTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("HPC Cache Blob Target not found: %s", resourceName)
		}

		id, err := parse.StorageTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).HPCCache.StorageTargetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HPC Cache Blob Target %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
		}

		return nil
	}
}

func testCheckHPCCacheBlobTargetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).HPCCache.StorageTargetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hpc_cache_blob_target" {
			continue
		}

		id, err := parse.StorageTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccHPCCacheBlobTarget_basic(data acceptance.TestData) string {
	template := testAccHPCCacheBlobTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage1"
}
`, template, data.RandomString)
}

func testAccHPCCacheBlobTarget_namespace(data acceptance.TestData) string {
	template := testAccHPCCacheBlobTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage2"
}
`, template, data.RandomString)
}

func testAccHPCCacheBlobTarget_requiresImport(data acceptance.TestData) string {
	template := testAccHPCCacheBlobTarget_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "import" {
  name                 = azurerm_hpc_cache_blob_target.test.name
  resource_group_name  = azurerm_hpc_cache_blob_target.test.resource_group_name
  cache_name           = azurerm_hpc_cache_blob_target.test.cache_name
  storage_container_id = azurerm_hpc_cache_blob_target.test.storage_container_id
  namespace_path       = azurerm_hpc_cache_blob_target.test.namespace_path
}
`, template)
}

func testAccHPCCacheBlobTarget_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
}

data "azuread_service_principal" "test" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorgacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctest-strgctn-%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_role_assignment" "test_storage_account_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}
