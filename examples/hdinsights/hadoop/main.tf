resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
  tags     = "${var.tags}"
}

resource "azurerm_storage_account" "main" {
  name                     = "${var.prefix}stor"
  location                 = "${azurerm_resource_group.main.location}"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
  tags                     = "${var.tags}"
}

resource "azurerm_storage_container" "main" {
  name                  = "data"
  resource_group_name   = "${azurerm_resource_group.main.name}"
  storage_account_name  = "${azurerm_storage_account.main.name}"
  container_access_type = "private"
}

resource "azurerm_hdinsight_cluster" "main" {
  name                = "${var.prefix}-hadoop"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  tier                = "standard"
  tags                = "${var.tags}"

  cluster {
    kind    = "Hadoop"
    version = "3.6"

    gateway {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  storage_profile {
    storage_account {
      storage_account_name = "${azurerm_storage_account.main.primary_blob_domain}"
      storage_account_key  = "${azurerm_storage_account.main.primary_access_key}"
      container_name       = "${azurerm_storage_container.main.name}"
      is_default           = true
    }
  }

  head_node {
    target_instance_count = 2

    hardware_profile {
      vm_size = "Standard_D12_V2"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  worker_node {
    target_instance_count = 4

    hardware_profile {
      vm_size = "Standard_D4_V2"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  zookeeper_node {
    target_instance_count = 3

    hardware_profile {
      vm_size = "Medium"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }
}
