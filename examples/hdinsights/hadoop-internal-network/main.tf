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
}

resource "azurerm_storage_container" "main" {
  name                  = "data"
  resource_group_name   = "${azurerm_resource_group.main.name}"
  storage_account_name  = "${azurerm_storage_account.main.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_hdinsight_cluster" "main" {
  name                = "${var.prefix}-hadoop"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  tier                = "standard"

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

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.main.id}"
      subnet_id          = "${azurerm_subnet.internal.id}"
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

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.main.id}"
      subnet_id          = "${azurerm_subnet.internal.id}"
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

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.main.id}"
      subnet_id          = "${azurerm_subnet.internal.id}"
    }
  }
}
