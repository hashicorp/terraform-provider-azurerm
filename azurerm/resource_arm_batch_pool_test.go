package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBatchPool_basic(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPool_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPool_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
			{
				Config:      testaccAzureRMBatchPool_requiresImport(ri, rs, testLocation()),
				ExpectError: testRequiresImportError("azurerm_batch_account"),
			},
		},
	})
}

func TestAccAzureRMBatchPool_fixedScale_complete(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPool_fixedScale_complete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "max_tasks_per_node", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_autoScale_complete(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPool_autoScale_complete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.0.evaluation_interval", "PT15M"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_completeUpdated(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPool_fixedScale_complete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
			{
				Config: testaccAzureRMBatchPool_autoScale_complete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.0.evaluation_interval", "PT15M"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPoolStartTask_basic(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPoolStartTask_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(resourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(resourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(resourceName, "start_task.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.max_task_retry_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.environment.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.environment.env", "TEST"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.environment.bu", "Research&Dev"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.user_identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.user_identity.0.auto_user.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.user_identity.0.auto_user.0.scope", "Task"),
					resource.TestCheckResourceAttr(resourceName, "start_task.0.user_identity.0.auto_user.0.elevation_level", "NonAdmin"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_certificates(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificate0ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, ri, rs)
	certificate1ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, ri, rs)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testaccAzureRMBatchPoolCertificates(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.id", certificate0ID),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.store_location", "CurrentUser"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.store_name", ""),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.visibility.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.visibility.3294600504", "StartTask"),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.id", certificate1ID),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.store_location", "CurrentUser"),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.store_name", ""),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.visibility.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.visibility.3294600504", "StartTask"),
					resource.TestCheckResourceAttr(resourceName, "certificate.1.visibility.4077195354", "RemoteUser"),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_validateResourceFileWithoutSource(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testaccAzureRMBatchPoolValidateResourceFileWithoutSource(ri, rs, testLocation()),
				ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
			},
		},
	})
}

func TestAccAzureRMBatchPool_validateResourceFileWithMultipleSources(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testaccAzureRMBatchPoolValidateResourceFileWithMultipleSources(ri, rs, testLocation()),
				ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
			},
		},
	})
}

func TestAccAzureRMBatchPool_validateResourceFileBlobPrefixWithoutAutoStorageContainerUrl(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testaccAzureRMBatchPoolValidateResourceFileBlobPrefixWithoutAutoStorageContainerName(ri, rs, testLocation()),
				ExpectError: regexp.MustCompile("auto_storage_container_name or storage_container_url must be specified when using blob_prefix"),
			},
		},
	})
}

func TestAccAzureRMBatchPool_validateResourceFileHttpURLWithoutFilePath(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testaccAzureRMBatchPoolValidateResourceFileHttpURLWithoutFilePath(ri, rs, testLocation()),
				ExpectError: regexp.MustCompile("file_path must be specified when using http_url"),
			},
		},
	})
}

func testCheckAzureRMBatchPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		poolName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).batchPoolClient

		resp, err := conn.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on batchPoolClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Batch pool %q (account: %q, resource group: %q) does not exist", poolName, accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMBatchPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_pool" {
			continue
		}

		poolName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).batchPoolClient

		resp, err := conn.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testaccAzureRMBatchPool_fixedScale_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
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

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  max_tasks_per_node  = 2
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  
  fixed_scale {
    target_dedicated_nodes = 2
  }
  
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, rInt, location, rString, rString, rString)
}

func testaccAzureRMBatchPool_autoScale_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
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

resource "azurerm_batch_pool" "test" {
  name                          = "testaccpool%s"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  account_name                  = "${azurerm_batch_account.test.name}"
  display_name                  = "Test Acc Pool"
  vm_size                       = "Standard_A1"
  node_agent_sku_id             = "batch.node.ubuntu 16.04"
  stop_pending_resize_operation = true

  auto_scale {
    evaluation_interval = "PT15M"
    formula             = <<EOF
      startingNumberOfVMs = 1;
      maxNumberofVMs = 25;
      pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
      pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 * TimeInterval_Second));
      $TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
EOF
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, rInt, location, rString, rString, rString)
}

func testaccAzureRMBatchPool_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPool_requiresImport(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_batch_pool" "import" {
  name                = "${azurerm_batch_pool.test.name}"
  resource_group_name = "${azurerm_batch_pool.test.resource_group_name}"
  account_name        = "${azurerm_batch_pool.test.account_name}"
  node_agent_sku_id   = "${azurerm_batch_pool.test.node_agent_sku_id}"
  vm_size             = "${azurerm_batch_pool.test.vm_size}"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, testaccAzureRMBatchPool_basic(rInt, rString, location))
}

func testaccAzureRMBatchPoolStartTask_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST",
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url  = "https://raw.githubusercontent.com/terraform-providers/terraform-provider-azurerm/master/README.md"
      file_path = "README.md"
    }
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolValidateResourceFileWithoutSource(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST",
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      # no valid values for sources
      auto_storage_container_name = ""
      file_mode                   = "0770"
    }
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolValidateResourceFileWithMultipleSources(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST",
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      auto_storage_container_name = "test"
      http_url  = "test"
      file_path = "README.md"
    }
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolValidateResourceFileBlobPrefixWithoutAutoStorageContainerName(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST",
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url    = "test"
      blob_prefix = "test"
      file_path   = "README.md"
    }
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolValidateResourceFileHttpURLWithoutFilePath(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST",
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url  = "test"
      file_path = ""
    }
  }
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolCertificates(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_certificate" "testcer" {
	resource_group_name  = "${azurerm_resource_group.test.name}"
	account_name         = "${azurerm_batch_account.test.name}"
	certificate          = "${base64encode(file("testdata/batch_certificate.cer"))}"
	format               = "Cer"
	thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34" # deliberately using lowercase here as verification
	thumbprint_algorithm = "SHA1"
}
resource "azurerm_batch_certificate" "testpfx" {
	resource_group_name  = "${azurerm_resource_group.test.name}"
	account_name         = "${azurerm_batch_account.test.name}"
	certificate          = "${base64encode(file("testdata/batch_certificate.pfx"))}"
	format               = "Pfx"
	password             = "terraform"
	thumbprint           = "42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"
	thumbprint_algorithm = "SHA1"
}
	
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
	}
	certificate = [ {
		id             = "${azurerm_batch_certificate.testcer.id}"
		store_location = "CurrentUser"
		visibility     = [ "StartTask" ]
	}, {
		id             = "${azurerm_batch_certificate.testpfx.id}"
		store_location = "CurrentUser"
		visibility     = [ "StartTask", "RemoteUser" ]
	}]
}

`, rInt, location, rString, rString)
}
