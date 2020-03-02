package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMHPCCacheBlobTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHPCCacheBlobTarget_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHPCCacheBlobTarget_namespace(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHPCCacheBlobTarget_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheBlobTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheBlobTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheBlobTargetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMHPCCacheBlobTarget_requiresImport),
		},
	})
}

func testCheckAzureRMHPCCacheBlobTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("HPC Cache Blob Target not found: %s", resourceName)
		}

		id, err := parsers.HPCCacheTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.StorageTargetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Cache, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HPC Cache Blob Target %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMHPCCacheBlobTargetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.StorageTargetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hpc_cache_blob_target" {
			continue
		}

		id, err := parsers.HPCCacheTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Cache, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMHPCCacheBlobTarget_basic(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheBlobTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.id
  namespace_path       = "/blob_storage1"
}
`, template, data.RandomString)
}

func testAccAzureRMHPCCacheBlobTarget_namespace(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheBlobTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.id
  namespace_path       = "/blob_storage2"
}
`, template, data.RandomString)
}

func testAccAzureRMHPCCacheBlobTarget_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheBlobTarget_basic(data)
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

func testAccAzureRMHPCCacheBlobTarget_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azuread_service_principal" "test" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorgacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctest-strgctn-%[2]s"
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
`, testAccAzureRMHPCCache_basic(data), data.RandomString)
}
