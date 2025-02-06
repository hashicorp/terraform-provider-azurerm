// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BatchPoolDataSource struct{}

func TestAccBatchPoolDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_pool", "test")
	r := BatchPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("microsoft-azure-batch"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("20-04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("ubuntu-server-container"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 20.04"),
				check.That(data.ResourceName).Key("max_tasks_per_node").HasValue("2"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.task_retry_maximum").HasValue("5"),
				check.That(data.ResourceName).Key("start_task.0.common_environment_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.common_environment_properties.env").HasValue("TEST"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.scope").HasValue("Task"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.elevation_level").HasValue("NonAdmin"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.id").Exists(),
				check.That(data.ResourceName).Key("certificate.0.store_location").HasValue("CurrentUser"),
				check.That(data.ResourceName).Key("certificate.0.store_name").HasValue(""),
				check.That(data.ResourceName).Key("certificate.0.visibility.#").HasValue("2"),
				check.That(data.ResourceName).Key("container_configuration.0.type").HasValue("DockerCompatible"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.registry_server").HasValue("myContainerRegistry.azurecr.io"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("metadata.tagName").HasValue("Example tag"),
			),
		},
	})
}

func (BatchPoolDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "acctest-publicip-%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                                = "testaccbatch%s"
  resource_group_name                 = azurerm_resource_group.test.name
  location                            = azurerm_resource_group.test.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.test.id
  storage_account_authentication_mode = "StorageKeys"

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate_password.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 20.04"
  max_tasks_per_node  = 2

  fixed_scale {
    target_dedicated_nodes = 2
    resize_timeout         = "PT15M"
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
    version   = "latest"
  }

  certificate {
    id             = azurerm_batch_certificate.test.id
    store_location = "CurrentUser"
    visibility     = ["StartTask", "RemoteUser"]
  }

  container_configuration {
    type = "DockerCompatible"
    container_registries {
      registry_server = "myContainerRegistry.azurecr.io"
      user_name       = "myUserName"
      password        = "myPassword"
    }
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 5
    wait_for_success   = true

    common_environment_properties = {
      env = "TEST"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }
  }

  network_configuration {
    public_address_provisioning_type = "UserManaged"
    public_ips                       = [azurerm_public_ip.test.id]
    subnet_id                        = azurerm_subnet.test.id
  }

  metadata = {
    tagName = "Example tag"
  }
}

data "azurerm_batch_pool" "test" {
  name                = azurerm_batch_pool.test.name
  account_name        = azurerm_batch_pool.test.account_name
  resource_group_name = azurerm_batch_pool.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString)
}
