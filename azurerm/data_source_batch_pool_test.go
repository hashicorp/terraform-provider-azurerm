package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMBatchPool_complete(t *testing.T) {
	dataSourceName := "data.azurerm_batch_pool.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMBatchPool_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccpool%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.max_task_retry_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.environment.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.environment.env", "TEST"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.0.scope", "Task"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.0.elevation_level", "NonAdmin"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMBatchPool_complete(rInt int, rString string, location string) string {
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
  name                   = "testaccpool%s"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  account_name           = "${azurerm_batch_account.test.name}"
  display_name           = "Test Acc Pool"
  vm_size                = "Standard_A1"
  node_agent_sku_id      = "batch.node.ubuntu 16.04"

  fixed_scale {
    target_dedicated_nodes = 2
    resize_timeout         = "PT15M"
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
      env = "TEST"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }
  }
}

data "azurerm_batch_pool" "test" {
  name                = "${azurerm_batch_pool.test.name}"
  account_name        = "${azurerm_batch_pool.test.account_name}"
  resource_group_name = "${azurerm_batch_pool.test.resource_group_name}"
}
`, rInt, location, rString, rString, rString)
}
