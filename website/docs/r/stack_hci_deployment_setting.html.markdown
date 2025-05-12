---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_deployment_setting"
description: |-
  Manages a Stack HCI Deployment Setting.
---

# azurerm_stack_hci_deployment_setting

Manages a Stack HCI Deployment Setting.

-> **Note:** Completion of the prerequisites of deploying the Azure Stack HCI in your environment is outside the scope of this document. For more details refer to the [Azure Stack HCI deployment sequence](https://learn.microsoft.com/en-us/azure-stack/hci/deploy/deployment-introduction#deployment-sequence). If you encounter issues completing the prerequisites, we'd recommend opening a ticket with Microsoft Support.

-> **Note:** During the deployment process, the service will generate additional resources, including a new Arc Bridge Appliance and a Custom Location containing several Stack HCI Storage Paths. The provider will attempt to remove these resources on the deletion or recreation of `azurerm_stack_hci_deployment_setting`.

## Example Usage

```hcl
provider "azuread" {}

variable "local_admin_user" {
  description = "The username of the local administrator account."
  sensitive   = true
  type        = string
}

variable "local_admin_password" {
  description = "The password of the local administrator account."
  sensitive   = true
  type        = string
}

variable "domain_admin_user" {
  description = "The username of the domain account."
  sensitive   = true
  type        = string
}

variable "domain_admin_password" {
  description = "The password of the domain account."
  sensitive   = true
  type        = string
}

variable "deployment_user" {
  sensitive   = true
  type        = string
  description = "The username for deployment user."
}

variable "deployment_user_password" {
  sensitive   = true
  type        = string
  description = "The password for deployment user."
}

locals {
  servers = [
    {
      name        = "AzSHOST1",
      ipv4Address = "192.168.1.12"
    },
    {
      name        = "AzSHOST2",
      ipv4Address = "192.168.1.13"
    }
  ]
  connection_roles = [
    "Azure Connected Machine Onboarding",
    "Azure Connected Machine Resource Administrator",
    "Azure Resource Bridge Deployment Role"
  ]
  machine_roles = [
    "Key Vault Secrets User",
    "Azure Connected Machine Resource Manager",
    "Azure Stack HCI Device Management Role",
    "Reader"
  ]
}

data "azurerm_resource_group" "example" {
  name = "hci-example"
}

resource "azuread_application" "example" {
  display_name = "example-hci-onboard"
}

# https://learn.microsoft.com/en-us/azure-stack/hci/deploy/deployment-azure-resource-manager-template#create-a-service-principal-and-client-secret
resource "azuread_service_principal" "example" {
  client_id = azuread_application.example.client_id
}

resource "azuread_service_principal_password" "example" {
  service_principal_id = azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "connect" {
  count                = length(local.connection_roles)
  scope                = data.azurerm_resource_group.example.id
  role_definition_name = local.connection_roles[count.index]
  principal_id         = azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "connect" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Contributor"
  principal_id         = azuread_service_principal.example.object_id
}

# after preparing Active Directory and registering servers with Arc, and then the Arc VM is ready
data "azurerm_arc_machine" "server" {
  count               = length(local.servers)
  name                = local.servers[count.index].name
  resource_group_name = data.azurerm_resource_group.example.name
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "DeploymentKeyVault" {
  name                            = "hci-examplekv"
  location                        = data.azurerm_resource_group.example.location
  resource_group_name             = data.azurerm_resource_group.example.name
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  enabled_for_disk_encryption     = true
  tenant_id                       = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days      = 30
  enable_rbac_authorization       = true
  public_network_access_enabled   = true
  sku_name                        = "standard"
}

resource "azurerm_role_assignment" "KeyVault" {
  scope                = azurerm_key_vault.DeploymentKeyVault.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_secret" "AzureStackLCMUserCredential" {
  name         = "AzureStackLCMUserCredential"
  content_type = "Secret"
  value        = base64encode("${var.deployment_user}:${var.deployment_user_password}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "LocalAdminCredential" {
  name         = "LocalAdminCredential"
  content_type = "Secret"
  value        = base64encode("${var.local_admin_user}:${var.local_admin_password}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "DefaultARBApplication" {
  name         = "DefaultARBApplication"
  content_type = "Secret"
  value        = base64encode("${azuread_service_principal.example.client_id}:${azuread_service_principal_password.example.value}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "WitnessStorageKey" {
  name         = "WitnessStorageKey"
  content_type = "Secret"
  value        = base64encode(azurerm_storage_account.witness.primary_access_key)
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_storage_account" "witness" {
  name                     = "hciexamplesta"
  location                 = data.azurerm_resource_group.example.location
  resource_group_name      = data.azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

// service principal of 'Microsoft.AzureStackHCI Resource Provider' application
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "MachineRoleAssign1" {
  count                = length(local.machine_roles)
  scope                = data.azurerm_resource_group.example.id
  role_definition_name = local.machine_roles[count.index]
  principal_id         = data.azurerm_arc_machine.server[0].identity[0].principal_id
}

resource "azurerm_role_assignment" "MachineRoleAssign2" {
  count                = length(local.machine_roles)
  scope                = data.azurerm_resource_group.example.id
  role_definition_name = local.machine_roles[count.index]
  principal_id         = data.azurerm_arc_machine.server[1].identity[0].principal_id
}

resource "azurerm_role_assignment" "ServicePrincipalRoleAssign" {
  scope                = data.azurerm_resource_group.example.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}

resource "azurerm_stack_hci_cluster" "example" {
  depends_on = [
    azurerm_key_vault_secret.DefaultARBApplication,
    azurerm_key_vault_secret.AzureStackLCMUserCredential,
    azurerm_key_vault_secret.LocalAdminCredential,
    azurerm_key_vault_secret.WitnessStorageKey,
    azurerm_role_assignment.connect,
    azurerm_role_assignment.connect2,
    azurerm_role_assignment.ServicePrincipalRoleAssign,
    azurerm_role_assignment.MachineRoleAssign1,
    azurerm_role_assignment.MachineRoleAssign2,
  ]
  name                = "hci-cl"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  identity {
    type = "SystemAssigned"
  }

  // the client_id will be populated after deployment
  lifecycle {
    ignore_changes = [
      client_id
    ]
  }
}

resource "azurerm_stack_hci_deployment_setting" "example" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.example.id
  arc_resource_ids     = [for server in data.azurerm_arc_machine.server : server.id]
  version              = "10.0.0.0"

  scale_unit {
    active_directory_organizational_unit_path = "OU=hci,DC=jumpstart,DC=local"
    domain_fqdn                               = "jumpstart.local"
    secrets_location                          = azurerm_key_vault.DeploymentKeyVault.vault_uri
    name_prefix                               = "hci"

    cluster {
      azure_service_endpoint = "core.windows.net"
      cloud_account_name     = azurerm_storage_account.witness.name
      name                   = azurerm_stack_hci_cluster.example.name
      witness_type           = "Cloud"
      witness_path           = "Cloud"
    }

    host_network {
      storage_auto_ip_enabled                 = true
      storage_connectivity_switchless_enabled = false
      intent {
        name                                          = "ManagementCompute"
        adapter_property_override_enabled             = false
        qos_policy_override_enabled                   = false
        virtual_switch_configuration_override_enabled = false
        adapter = [
          "FABRIC",
          "FABRIC2",
        ]
        traffic_type = [
          "Management",
          "Compute",
        ]
        qos_policy_override {
          priority_value8021_action_cluster = "7"
          priority_value8021_action_smb     = "3"
          bandwidth_percentage_smb          = "50"
        }
        adapter_property_override {
          jumbo_packet              = "9014"
          network_direct            = "Disabled"
          network_direct_technology = "RoCEv2"
        }
      }

      intent {
        name                                          = "Storage"
        adapter_property_override_enabled             = false
        qos_policy_override_enabled                   = false
        virtual_switch_configuration_override_enabled = false
        adapter = [
          "StorageA",
          "StorageB",
        ]
        traffic_type = [
          "Storage",
        ]
        qos_policy_override {
          priority_value8021_action_cluster = "7"
          priority_value8021_action_smb     = "3"
          bandwidth_percentage_smb          = "50"
        }
        adapter_property_override {
          jumbo_packet              = "9014"
          network_direct            = "Enabled"
          network_direct_technology = "RoCEv2"
        }
      }

      storage_network {
        name                 = "Storage1Network"
        network_adapter_name = "StorageA"
        vlan_id              = "711"
      }

      storage_network {
        name                 = "Storage2Network"
        network_adapter_name = "StorageB"
        vlan_id              = "712"
      }
    }

    infrastructure_network {
      gateway      = "192.168.1.1"
      subnet_mask  = "255.255.255.0"
      dhcp_enabled = false
      dns_server = [
        "192.168.1.254"
      ]
      ip_pool {
        ending_address   = "192.168.1.65"
        starting_address = "192.168.1.55"
      }
    }

    optional_service {
      custom_location = "customlocation"
    }

    physical_node {
      ipv4_address = "192.168.1.12"
      name         = "AzSHOST1"
    }

    physical_node {
      ipv4_address = "192.168.1.13"
      name         = "AzSHOST2"
    }

    observability {
      streaming_data_client_enabled = true
      eu_location_enabled           = false
      episodic_data_upload_enabled  = true
    }

    security_setting {
      bitlocker_boot_volume_enabled   = true
      bitlocker_data_volume_enabled   = true
      credential_guard_enabled        = true
      drift_control_enabled           = true
      drtm_protection_enabled         = true
      hvci_protection_enabled         = true
      side_channel_mitigation_enabled = true
      smb_cluster_encryption_enabled  = false
      smb_signing_enabled             = true
      wdac_enabled                    = true
    }

    storage {
      configuration_mode = "Express"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `stack_hci_cluster_id` - (Required) The ID of the Azure Stack HCI cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `arc_resource_ids` - (Required) Specifies a list of IDs of Azure ARC machine resource to be part of cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `scale_unit` - (Required) One or more `scale_unit` blocks as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `version` - (Required) The deployment template version. The format must be a set of numbers separated by dots such as `10.0.0.0`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `adapter_property_override` block supports the following:

* `jumbo_packet` - (Optional) The jumbo frame size of the adapter. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

* `network_direct` - (Optional) The network direct of the adapter. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

* `network_direct_technology` - (Optional) The network direct technology of the adapter. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `cluster` block supports the following:

* `azure_service_endpoint` - (Required) Specifies the Azure blob service endpoint, for example, `core.windows.net`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `cloud_account_name` - (Required) Specifies the Azure Storage account name of the cloud witness for the Azure Stack HCI cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `name` - (Required) Specifies the name of the cluster. It must be 3-15 characters long and contain only letters, numbers and hyphens. Changing this forces a new Stack HCI Deployment Setting to be created.

* `witness_path` - (Required) Specifies the fileshare path of the local witness for the Azure Stack HCI cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `witness_type` - (Required) Specifies the type of the witness. Possible values are `Cloud`, `FileShare`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `host_network` block supports the following:

* `intent` - (Required) One or more `intent` blocks as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `storage_network` - (Required) One or more `storage_network` blocks as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `storage_auto_ip_enabled` - (Optional) Whether allows users to specify IPs and Mask for Storage NICs when Network ATC is not assigning the IPs for storage automatically. Optional parameter required only for [3 nodes switchless deployments](https://learn.microsoft.com/azure-stack/hci/concepts/physical-network-requirements?tabs=overview%2C23H2reqs#using-switchless). Possible values are `true` and `false`. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `storage_connectivity_switchless_enabled` - (Optional) Defines how the storage adapters between nodes are connected either switch or switch less. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `infrastructure_network` block supports the following:

* `dns_server` - (Required) Specifies a list of IPv4 addresses of the DNS servers in your environment. Changing this forces a new Stack HCI Deployment Setting to be created.

* `gateway` - (Required) Specifies the default gateway that should be used for the provided IP address space. It should be in the format of an IPv4 IP address. Changing this forces a new Stack HCI Deployment Setting to be created.

* `ip_pool` - (Required) One or more `ip_pool` blocks as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `subnet_mask` - (Required) Specifies the subnet mask that matches the provided IP address space. Changing this forces a new Stack HCI Deployment Setting to be created.

* `dhcp_enabled` - (Optional) Whether DHCP is enabled for hosts and cluster IPs. Possible values are `true` and `false`. defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

-> **Note:** If `dhcp_enabled` is set to `false`, the deployment will use static IPs. If set to `true`, the gateway and DNS servers are not required.

---

A `intent` block supports the following:

* `name` - (Required) Specifies the name of the intent. Changing this forces a new Stack HCI Deployment Setting to be created.

* `adapter` - (Required) Specifies a list of ID of network interfaces used for the network intent. Changing this forces a new Stack HCI Deployment Setting to be created.

* `traffic_type` - (Required) Specifies a list of network traffic types. Possible values are `Compute`, `Storage`, `Management`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `adapter_property_override` - (Optional) A `adapter_property_override` block as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `adapter_property_override_enabled` - (Optional) Whether to override adapter properties. Possible values are `true` and `false`. defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `qos_policy_override` - (Optional) A `qos_policy_override` block as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `qos_policy_override_enabled` - (Optional) Whether to override QoS policy. Possible values are `true` and `false`. defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `virtual_switch_configuration_override` - (Optional) A `virtual_switch_configuration_override` block as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `virtual_switch_configuration_override_enabled` - (Optional) Whether to override virtual switch configuration. Possible values are `true` and `false`. defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `ip_pool` block supports the following:

* `ending_address` - (Required) Specifies starting IP address for the management network. A minimum of six free, contiguous IPv4 addresses (excluding your host IPs) are needed for infrastructure services such as clustering. Changing this forces a new Stack HCI Deployment Setting to be created.

* `starting_address` - (Required) Specifies ending IP address for the management network. A minimum of six free, contiguous IPv4 addresses (excluding your host IPs) are needed for infrastructure services such as clustering. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `observability` block supports the following:

* `episodic_data_upload_enabled` - (Required) Whether to collect log data to facilitate quicker issue resolution. Possible values are `true` and `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `eu_location_enabled` - (Required) Whether to store log and diagnostic data sent to Microsoft in EU or outside of EU. The log and diagnostic data is sent to the appropriate diagnostics servers depending upon where your cluster resides. Setting this to `false` results in all data sent to Microsoft to be stored outside of EU. Possible values are `true` and `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `streaming_data_client_enabled` - (Required) Whether the telemetry data will be sent to Microsoft. Possible values are `true` and `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `optional_service` block supports the following:

* `custom_location` - (Required) Specifies the name of custom location. A custom location will be created after the deployment is completed. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `physical_node` block supports the following:

* `ipv4_address` - (Required) Specifies the IPv4 address assigned to each physical server on your Azure Stack HCI cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `name` - (Required) The NETBIOS name of each physical server on your Azure Stack HCI cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `qos_policy_override` block supports the following:

* `bandwidth_percentage_smb` - (Optional) Specifies the percentage of the allocated storage traffic bandwidth. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

* `priority_value8021_action_cluster` - (Optional) Specifies the Cluster traffic priority. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

* `priority_value8021_action_smb` - (Optional) Specifies the Priority Flow Control where Data Center Bridging (DCB) is used. This parameter should only be modified based on your OEM guidance. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `scale_unit` block supports the following:

* `active_directory_organizational_unit_path` - (Required) Specify the full name of the Active Directory Organizational Unit container object prepared for the deployment, including the domain components. For example:`OU=HCI01,DC=contoso,DC=com`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `cluster` - (Required) A `cluster` block as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `domain_fqdn` - (Required) Specifies the FQDN for deploying cluster. Changing this forces a new Stack HCI Deployment Setting to be created.

* `host_network` - (Required) A `host_network` block as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `infrastructure_network` - (Required) One or more `infrastructure_network` blocks as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `name_prefix` - (Required) Specifies the name prefix to deploy cluster. It must be 1-8 characters long and contain only letters, numbers and hyphens Changing this forces a new Stack HCI Deployment Setting to be created.

* `optional_service` - (Required) A `optional_service` block as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `physical_node` - (Required) One or more `physical_node` blocks as defined above. Changing this forces a new Stack HCI Deployment Setting to be created.

* `secrets_location` - (Required) The URI to the Key Vault or secret store. Changing this forces a new Stack HCI Deployment Setting to be created.

* `security_setting` - (Required) A `security_setting` block as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `storage` - (Required) A `storage` block as defined below. Changing this forces a new Stack HCI Deployment Setting to be created.

* `episodic_data_upload_enabled` - (Optional) Whether to collect log data to facilitate quicker issue resolution. Possible values are `true` and `false`. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `eu_location_enabled` - (Optional) Whether to store data sent to Microsoft in EU. The log and diagnostic data is sent to the appropriate diagnostics servers depending upon where your cluster resides. Setting this to `false` results in all data sent to Microsoft to be stored outside of the EU. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `streaming_data_client_enabled` - (Optional) Whether the telemetry data will be sent to Microsoft. Possible values are `true` and `false`. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `bitlocker_boot_volume_enabled` - (Optional) Whether to enable BitLocker for boot volume. Possible values are `true` and `false`. When set to `true`, BitLocker XTS_AES 256-bit encryption is enabled for all data-at-rest on the OS volume of your Azure Stack HCI cluster. This setting is TPM-hardware dependent. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `bitlocker_data_volume_enabled` - (Optional) Whether to enable BitLocker for data volume. Possible values are `true` and `false`. When set to `true`, BitLocker XTS-AES 256-bit encryption is enabled for all data-at-rest on your Azure Stack HCI cluster shared volumes. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `credential_guard_enabled` - (Optional) Whether to enable credential guard. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `drift_control_enabled` - (Optional) Whether to enable drift control. Possible values are `true` and `false`. When set to `true`, the security baseline is re-applied regularly. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `drtm_protection_enabled` - (Optional) Whether to enable DRTM protection. Possible values are `true` and `false`. When set to `true`, Secure Boot is enabled on your Azure HCI cluster. This setting is hardware dependent. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `hvci_protection_enabled` - (Optional) Whether to enable HVCI protection. Possible values are `true` and `false`. When set to `true`, Hypervisor-protected Code Integrity is enabled on your Azure HCI cluster. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `side_channel_mitigation_enabled` - (Optional) Whether to enable side channel mitigation. Possible values are `true` and `false`. When set to `true`, all side channel mitigations are enabled on your Azure HCI cluster. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `smb_cluster_encryption_enabled` - (Optional) Whether to enable SMB cluster encryption. Possible values are `true` and `false`. When set to `true`, cluster east-west traffic is encrypted. Defaults to `false`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `smb_signing_enabled` - (Optional) Whether to enable SMB signing. Possible values are `true` and `false`. When set to `true`, the SMB default instance requires sign in for the client and server services. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

* `wdac_enabled` - (Optional) Whether to enable WDAC. Possible values are `true` and `false`. When set to `true`, applications and the code that you can run on your Azure Stack HCI cluster are limited. Defaults to `true`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `storage` block supports the following:

* `configuration_mode` - (Required) The configuration mode of storage. If set to `Express` and your storage is configured as per best practices based on the number of nodes in the cluster. Possible values are `Express`, `InfraOnly` and `KeepStorage`. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `storage_network` block supports the following:

* `name` - (Required) The name of the storage network. Changing this forces a new Stack HCI Deployment Setting to be created.

* `network_adapter_name` - (Required) The name of the network adapter. Changing this forces a new Stack HCI Deployment Setting to be created.

* `vlan_id` - (Required) Specifies the ID for the VLAN storage network. This setting is applied to the network interfaces that route the storage and VM migration traffic. Changing this forces a new Stack HCI Deployment Setting to be created.

---

A `virtual_switch_configuration_override` block supports the following:

* `enable_iov` - (Optional) Specifies the IoV enable status for Virtual Switch. Changing this forces a new Stack HCI Deployment Setting to be created.

* `load_balancing_algorithm` - (Optional) Specifies the load balancing algorithm for Virtual Switch. Changing this forces a new Stack HCI Deployment Setting to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Stack HCI Deployment Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Stack HCI Deployment Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stack HCI Deployment Setting.
* `delete` - (Defaults to 1 hour) Used when deleting the Stack HCI Deployment Setting.

## Import

Stack HCI Deployment Settings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_deployment_setting.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/clus1/deploymentSettings/default
```
