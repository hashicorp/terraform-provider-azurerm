# Configure the Microsoft Azure Provider
provider "azurerm" {
  # if you're using a Service Principal (shared account) then either set the environment variables, or fill these in:   # subscription_id = "..."  # client_id       = "..."   # client_secret   = "..."  # tenant_id       = "..."
}

resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_storage_account" "stor" {
  name                     = "stor${random_integer.ri.result}"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  location                 = "${azurerm_resource_group.rg.location}"
  account_tier             = "${var.storage_account_tier}"
  account_replication_type = "${var.storage_replication_type}"
}

resource "azurerm_batch_account" "batch" {
  name                 = "batch${random_integer.ri.result}"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  location             = "${azurerm_resource_group.rg.location}"
  storage_account_name = "${azurerm_storage_account.stor.name}"
}

resource "azurerm_batch_pool" "fixedpool" {
  name                = "myfixedpool"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  account_name        = "${azurerm_batch_account.batch.name}"
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
}

resource "azurerm_batch_pool" "autopool" {
  name                = "myautopool"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  account_name        = "${azurerm_batch_account.batch.name}"
  display_name        = "Auto Scale Pool"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  
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
