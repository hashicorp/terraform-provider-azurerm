package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAzureRMBatchPool_complete(t *testing.T) {
	dataSourceName := "data.azurerm_batch_pool.test"
	ri := tf.AccRandTimeInt()

	rs := acctest.RandString(4)
	location := acceptance.Location()
	config := testAccDataSourceAzureRMBatchPool_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccpool%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.publisher", "microsoft-azure-batch"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.sku", "16-04-lts"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_image_reference.0.offer", "ubuntu-server-container"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(dataSourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(dataSourceName, "max_tasks_per_node", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.max_task_retry_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.environment.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.environment.env", "TEST"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.0.scope", "Task"),
					resource.TestCheckResourceAttr(dataSourceName, "start_task.0.user_identity.0.auto_user.0.elevation_level", "NonAdmin"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.0.store_location", "CurrentUser"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.0.store_name", ""),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.0.visibility.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.0.visibility.3294600504", "StartTask"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate.0.visibility.4077195354", "RemoteUser"),
					resource.TestCheckResourceAttr(dataSourceName, "container_configuration.0.type", "DockerCompatible"),
					resource.TestCheckResourceAttr(dataSourceName, "container_configuration.0.container_registries.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "container_configuration.0.container_registries.0.registry_server", "myContainerRegistry.azurecr.io"),
					resource.TestCheckResourceAttr(dataSourceName, "container_configuration.0.container_registries.0.user_name", "myUserName"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMBatchPool_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batch"
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

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("testdata/batch_certificate.pfx")}"
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  max_tasks_per_node  = 2

  fixed_scale {
    target_dedicated_nodes = 2
    resize_timeout         = "PT15M"
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "16-04-lts"
    version   = "latest"
  }

  certificate {
    id             = "${azurerm_batch_certificate.test.id}"
    store_location = "CurrentUser"
    visibility     = ["StartTask", "RemoteUser"]
  }

  container_configuration {
    type = "DockerCompatible"
    container_registries = [
      {
        registry_server = "myContainerRegistry.azurecr.io"
        user_name       = "myUserName"
        password        = "myPassword"
      },
    ]
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
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
