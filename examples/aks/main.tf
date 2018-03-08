resource "azurerm_resource_group" "test" {
  name     = "${var.azurerm_resource_group_name}"
  location = "${var.azurerm_resource_group_location}"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "${var.azurerm_kubernetes_cluster_dns_prefix}"
  kubernetes_version  = "${var.azurerm_kubernetes_cluster_version}"

  linux_profile {
    admin_username = "${var.azurerm_kubernetes_cluster_admin_username}"

    ssh_key {
      key_data = "${var.azurerm_kubernetes_cluster_ssh_key}"
    }
  }

  agent_pool_profile {
    name    = "default"
    count   = "${var.azurerm_kubernetes_cluster_agent_count}"
    vm_size = "${var.azurerm_kubernetes_cluster_agent_vm_size}"
  }

  service_principal {
    client_id     = "${var.azurerm_kubernetes_cluster_client_id}"
    client_secret = "${var.azurerm_kubernetes_cluster_client_secret}"
  }
}