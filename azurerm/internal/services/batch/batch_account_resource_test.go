package batch_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateBatchAccountName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"ab", true},
		{"ABC", true},
		{"abc", false},
		{"123456789012345678901234", false},
		{"1234567890123456789012345", true},
		{"abc12345", false},
	}

	for _, test := range testCases {
		_, es := batch.ValidateAzureRMBatchAccountName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}

		if !test.shouldError && len(es) > 1 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccBatchAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "BatchService"),
				),
			},
		},
	})
}

func TestAccBatchAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchAccountExists(data.ResourceName),
				),
			},
			{
				Config:      testAccBatchAccount_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_batch_account"),
			},
		},
	})
}

func TestAccBatchAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccBatchAccount_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.version", "2"),
				),
			},
		},
	})
}

func TestAccBatchAccount_userSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_account", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchAccount_userSubscription(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "UserSubscription"),
				),
			},
		},
	})
}

func testCheckBatchAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.AccountClient

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Ensure resource group exists in API

		resp, err := conn.Get(ctx, id.ResourceGroup, id.BatchAccountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on batchAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Batch account %q (resource group: %q) does not exist", id.BatchAccountName, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckBatchAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_account" {
			continue
		}

		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.BatchAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccBatchAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccBatchAccount_requiresImport(data acceptance.TestData) string {
	template := testAccBatchAccount_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "import" {
  name                 = azurerm_batch_account.test.name
  resource_group_name  = azurerm_batch_account.test.resource_group_name
  location             = azurerm_batch_account.test.location
  pool_allocation_mode = azurerm_batch_account.test.pool_allocation_mode
}
`, template)
}

func testAccBatchAccount_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchAccount_completeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s2"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env     = "test"
    version = "2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchAccount_userSubscription(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azuread_service_principal" "test" {
  display_name = "Microsoft Azure Batch"
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                            = "batchkv%s"
  location                        = "${azurerm_resource_group.test.location}"
  resource_group_name             = "${azurerm_resource_group.test.name}"
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = "%s"

  sku_name = "standard"

  access_policy {
    tenant_id = "%s"
    object_id = "${data.azuread_service_principal.test.object_id}"

    secret_permissions = [
      "get",
      "list",
      "set",
      "delete"
    ]

  }
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  pool_allocation_mode = "UserSubscription"

  key_vault_reference {
    id  = "${azurerm_key_vault.test.id}"
    url = "${azurerm_key_vault.test.vault_uri}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tenantID, tenantID, data.RandomString)
}
