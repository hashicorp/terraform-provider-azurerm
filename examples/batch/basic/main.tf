provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stor"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "example" {
  name                = "${var.prefix}batch"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  storage_account_id  = "${azurerm_storage_account.example.id}"
}

resource "azurerm_batch_pool" "fixed" {
  name                = "${var.prefix}-fixed-pool"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_batch_account.example.name}"
  display_name        = "Fixed Scale Pool"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

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

  metadata ={
    "tagName"= "Example tag"
  }
}

resource "azurerm_batch_pool" "autopool" {
  name                = "${var.prefix}-autoscale-pool"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_batch_account.example.name}"
  display_name        = "Auto Scale Pool"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

  auto_scale {
    evaluation_interval = "PT15M"

    formula = <<EOF
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
