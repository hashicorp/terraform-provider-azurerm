package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
		_, es := validateAzureRMBatchAccountName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}

		if !test.shouldError && len(es) > 1 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccAzureRMBatchAccount_basic(t *testing.T) {
	resourceName := "azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	config := testAccAzureRMBatchAccount_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "pool_allocation_mode", "BatchService"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()

	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBatchAccount_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchAccountExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMBatchAccount_requiresImport(ri, rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_batch_account"),
			},
		},
	})
}

func TestAccAzureRMBatchAccount_complete(t *testing.T) {
	resourceName := "azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	config := testAccAzureRMBatchAccount_complete(ri, rs, location)
	configUpdate := testAccAzureRMBatchAccount_completeUpdated(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.version", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchAccount_userSubscription(t *testing.T) {
	resourceName := "azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	tenantID := os.Getenv("ARM_TENANT_ID")

	config := testAccAzureRMBatchAccount_userSubscription(ri, rs, location, tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "pool_allocation_mode", "UserSubscription"),
				),
			},
		},
	})
}

func testCheckAzureRMBatchAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		batchAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.AccountClient

		resp, err := conn.Get(ctx, resourceGroup, batchAccount)
		if err != nil {
			return fmt.Errorf("Bad: Get on batchAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Batch account %q (resource group: %q) does not exist", batchAccount, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMBatchAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMBatchAccount_basic(rInt int, batchAccountSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}
`, rInt, location, batchAccountSuffix)
}

func testAccAzureRMBatchAccount_requiresImport(rInt int, batchAccountSuffix string, location string) string {
	template := testAccAzureRMBatchAccount_basic(rInt, batchAccountSuffix, location)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "import" {
  name                 = "${azurerm_batch_account.test.name}"
  resource_group_name  = "${azurerm_batch_account.test.resource_group_name}"
  location             = "${azurerm_batch_account.test.location}"
  pool_allocation_mode = "${azurerm_batch_account.test.pool_allocation_mode}"
}
`, template)
}

func testAccAzureRMBatchAccount_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}
`, rInt, location, rString, rString)
}

func testAccAzureRMBatchAccount_completeUpdated(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s2"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env     = "test"
    version = "2"
  }
}
`, rInt, location, rString, rString)
}

func testAccAzureRMBatchAccount_userSubscription(rInt int, batchAccountSuffix string, location string, tenantID string) string {
	return fmt.Sprintf(`
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

  sku {
    name = "standard"
  }

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
`, rInt, location, batchAccountSuffix, tenantID, tenantID, batchAccountSuffix)
}
