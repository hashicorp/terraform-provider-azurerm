---
subcategory: "Service Fabric Managed Clusters"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_managed_cluster"
description: |-
  Manages a Resource Group.
---

# azurerm_service_fabric_managed_cluster

Manages a Resource Group.

## Example Usage

```hcl
resource "azurerm_service_fabric_managed_cluster" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
  http_gateway_port   = 4567

  lb_rule {
    backend_port       = 38080
    frontend_port      = 80
    probe_protocol     = "http"
    probe_request_path = "/test"
    protocol           = "tcp"
  }

  network_security_rules {
    access                       = "allow"
    protocol                     = "tcp"
    destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
    destination_port_ranges      = ["80", "443", "8080", "8190"]
    direction                    = "outbound"
    name                         = "test1"
    priority                     = 1000
    source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
    source_port_ranges           = ["10000-40000"]
  }

  client_connection_port = 12345

  node_type {
    data_disk_size_gb      = 130
    name                   = "test1"
    primary                = true
    application_port_range = "30000-49000"
    ephemeral_port_range   = "10000-20000"

    vm_size            = "Standard_DS1_v2"
    vm_image_publisher = "MicrosoftWindowsServer"
    vm_image_sku       = "2019-Datacenter-with-Containers"
    vm_image_offer     = "WindowsServer"
    vm_image_version   = "latest"
    vm_instance_count  = 5
  }
}
```

## Arguments Reference

The following arguments are supported:

* `client_connection_port` - (Required) Port to use when connecting to the cluster.

* `http_gateway_port` - (Required) Port that should be used by the Service Fabric Explorer to visualize applications and cluster status.

* `lb_rule` - (Required) One or more `lb_rule` blocks as defined below.

* `location` - (Required) The Azure Region where the Resource Group should exist. Changing this forces a new Resource Group to be created.

* `name` - (Required) The name which should be used for this Resource Group. Changing this forces a new Resource Group to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Group should exist.

---

* `authentication` - (Optional) Controls how connections to the cluster are authenticated. A `authentication` block as defined below.

* `backup_service_enabled` - (Optional) If true, backup service is enabled.

* `custom_fabric_setting` - (Optional) One or more `custom_fabric_setting` blocks as defined below.

* `dns_name` - (Optional) Hostname for the cluster. If unset the cluster's name will be used..

* `dns_service_enabled` - (Optional) If true, DNS service is enabled.

* `network_security_rules` - (Optional) One or more `network_security_rules` blocks as defined below.

* `node_type` - (Optional) One or more `node_type` blocks as defined below.

* `password` - (Optional) Administrator password for the VMs that will be created as part of this cluster.

* `sku` - (Optional) SKU for this cluster.  Changing this forces a new resource to be created. Default is `Basic`, allowed values are either `Basic` or `Standard`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Group.

* `upgrade_wave` - (Optional) Upgrade wave for the fabric runtime. Default is `Wave0`, allowed value must be one of `Wave0`, `Wave1`, or `Wave2`.

* `username` - (Optional) Administrator password for the VMs that will be created as part of this cluster.

---

A `active_directory` block supports the following:

* `client_application_id` - (Required) The ID of the Client Application.

* `cluster_application_id` - (Required) The ID of the Cluster Application.

* `tenant_id` - (Required) The ID of the Tenant.

---

A `authentication` block supports the following:

* `active_directory` - (Optional) A `active_directory` block as defined above.

* `certificate` - (Optional) One or more `certificate` blocks as defined below.

---

A `certificate` block supports the following:

* `thumbprint` - (Required) The thumbprint of the certificate.

* `type` - (Required) The type of the certificate. Can be `AdminClient` or `ReadOnlyClient`.

* `common_name` - (Optional) The certificate's CN.

---

A `certificates` block supports the following:

* `store` - (Required) The certificate store on the Virtual Machine to which the certificate should be added.

* `url` - (Required) The URL of a certificate that has been uploaded to Key Vault as a secret

---

A `custom_fabric_setting` block supports the following:

* `parameter` - (Required) Parameter name.

* `section` - (Required) Section name.

* `value` - (Required) Parameter value.

---

A `lb_rule` block supports the following:

* `backend_port` - (Required) LB Backend port.

* `frontend_port` - (Required) LB Frontend port.

* `probe_protocol` - (Required) Protocol for the probe. Can be one of `tcp`, `udp`, `http`, or `https`.

* `probe_request_path` - (Required) Path for the probe to check, when probe protocol is set to `http`.

* `protocol` - (Required) The transport protocol used in this rule. Can be one of `tcp` or `udp`.

---

A `node_type` block supports the following:

* `application_port_range` - (Required) Sets the port range available for applications. Format is `<from_port>-<to_port>`, for example `10000-20000`.

* `data_disk_size_gb` - (Required) The size of the data disk in gigabytes..

* `ephemeral_port_range` - (Required) Sets the port range available for the OS. Format is `<from_port>-<to_port>`, for example `10000-20000`. There has to be at least 255 ports available and cannot overlap with `application_port_range`..

* `name` - (Required) The name which should be used for this node type.

* `vm_image_offer` - (Required) The offer type of the marketplace image cluster VMs will use.

* `vm_image_publisher` - (Required) The publisher of the marketplace image cluster VMs will use.

* `vm_image_sku` - (Required) The SKU of the marketplace image cluster VMs will use.

* `vm_image_version` - (Required) The version of the marketplace image cluster VMs will use.

* `vm_instance_count` - (Required) The number of instances this node type will launch.

* `vm_size` - (Required) The size of the instances in this node type.

* `capacities` - (Optional) Specifies a list of key/value pairs used to set capacity tags for this node type.

* `data_disk_type` - (Optional) The type of the disk to use for storing data. It can be one of `Premium_LRS`, `Standard_LRS`, or `StandardSSD_LRS`.

* `multiple_placement_groups_enabled` - (Optional) If set the node type can be composed of multiple placement groups.

* `placement_properties` - (Optional) Specifies a list of placement tags that can be used to indicate where services should run..

* `primary` - (Optional) If set to true, system services will run on this node type. Only one node type should be marked as primary. Primary node type cannot be deleted or changed once they're created.

* `stateless` - (Optional) If set to true, only stateless workloads can run on this node type.

* `vm_secrets` - (Optional) One or more `vm_secrets` blocks as defined below.

---

A `network_security_rules` block supports the following:

* `access` - (Required) Specifies whether network traffic is allowed or denied. Possible values are `allow` and `deny`.

* `destination_address_prefixes` - (Required) List of destination address prefixes.

* `destination_port_ranges` - (Required) List of destination ports or port ranges.

* `direction` - (Required) The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are `inbound` and `outbound`.

* `name` - (Required) The name of the security rule.

* `priority` - (Required) Specifies the priority of the rule. The value can be between 1000 and 3000. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.

* `protocol` - (Required) Network protocol this rule applies to. Possible values include `ah`, `esp`, `http`, `https`, `icmp` ,`tcp` or `udp`.

* `source_address_prefixes` - (Required) List of source address prefixes.

* `source_port_ranges` - (Required) List of source ports or port ranges.

* `description` - (Optional) A description for this rule.

---

A `vm_secrets` block supports the following:

* `certificates` - (Required) One or more `certificates` blocks as defined above.

* `vault_id` - (Required) The ID of the Vault that contain the certificates.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour and 30 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 1 hour and 30 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 1 hour and 30 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_managed_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ServiceFabric/managedClusters/clusterName1
```
