package azurerm

import (
	"fmt"
	"net/http"
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
	location := testLocation()

	config := testaccAzureRMBatchPool_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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

func TestAccAzureRMBatchPool_fixedScale_complete(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()

	rs := acctest.RandString(4)
	location := testLocation()

	config := testaccAzureRMBatchPool_fixedScale_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
		},
	})
}

func TestAccAzureRMBatchPool_autoScale_complete(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()

	rs := acctest.RandString(4)
	location := testLocation()

	config := testaccAzureRMBatchPool_autoScale_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	location := testLocation()

	config := testaccAzureRMBatchPool_fixedScale_complete(ri, rs, location)
	configUpdate := testaccAzureRMBatchPool_autoScale_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
				Config: configUpdate,
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
	location := testLocation()

	config := testaccAzureRMBatchPoolStartTask_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
  name     = "testaccbatch%d"
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

  tags {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
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
  name     = "testaccbatch%d"
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

  tags {
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
}
`, rInt, location, rString, rString)
}

func testaccAzureRMBatchPoolStartTask_basic(rInt int, rString string, location string) string {
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

    environment {
			env = "TEST",
			bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }
  }
}
`, rInt, location, rString, rString)
}
