resource "azurerm_resource_group" "resource_group" {
  name     = "${var.azure_resource_group_name}"
  location = "${var.azure_resource_group_location}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_storage_account" "storage_account" {
  name                     = "${var.azure_storage_account_name}"
  resource_group_name      = "${azurerm_resource_group.resource_group.name}"
  location                 = "${azurerm_resource_group.resource_group.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_hdinsight_cluster" "hdinsight" {
  name                = "${var.azure_hdinsight_cluster_name}"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${azurerm_resource_group.resource_group.location}"
  cluster_version     = "${var.azure_hdinsight_cluster_version}"
  os_type             = "Linux"
  tier                = "Standard"

  cluster_definition {
    kind = "${var.azure_hdinsight_cluster_type}"

    configurations {
      gateway {
        rest_auth_credential__is_enabled = true
        rest_auth_credential__username   = "${var.azure_hdinsight_cluster_rest_username}"
        rest_auth_credential__password   = "${var.azure_hdinsight_cluster_rest_password}"
      }
    }
  }

  compute_profile {
    roles = [
      {
        name                  = "headnode"
        target_instance_count = "${var.azure_hdinsight_cluster_headnode_instance_count}"

        hardware_profile {
          vm_size = "${var.azure_hdinsight_cluster_headnode_vmsize}"
        }

        os_profile {
          linux_operating_system_profile {
            username = "${var.azure_hdinsight_cluster_headnode_username}"
            password = "${var.azure_hdinsight_cluster_headnode_password}"
          }
        }
      },
      {
        name                  = "workernode"
        target_instance_count = "${var.azure_hdinsight_cluster_workernode_instance_count}"

        hardware_profile {
          vm_size = "${var.azure_hdinsight_cluster_workernode_vmsize}"
        }

        os_profile {
          linux_operating_system_profile {
            username = "${var.azure_hdinsight_cluster_workernode_username}"

            ssh_key {
              key_data = "${var.azure_hdinsight_cluster_workernode_sshkey}"
            }
          }
        }
      },
    ]
  }

  storage_profile {
    storage_accounts = [
      {
        name       = "${azurerm_storage_account.storage_account.primary_blob_endpoint}"
        is_default = true
        container  = "${azurerm_resource_group.resource_group.name}"
        key        = "${azurerm_storage_account.storage_account.primary_access_key}"
      },
    ]
  }
}
