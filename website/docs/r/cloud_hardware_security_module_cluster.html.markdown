---
subcategory: "Hardware Security Module"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cloud_hardware_security_module_cluster"
description: |-
  Manages a Cloud Hardware Security Module Cluster.
---

# azurerm_cloud_hardware_security_module_cluster

Manages a Cloud Hardware Security Module Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cloud_hardware_security_module_cluster" "example" {
  name                = "example-cloudhsm-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Hardware Security Module Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud Hardware Security Module Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Cloud Hardware Security Module Cluster should exist. Changing this forces a new resource to be created.

---

* `identity` - (Optional) An `identity` block as defined below.

* `auto_generated_domain_name_label_scope` - (Optional) Specifies the scope for auto-generated domain name labels. Possible values are `TenantReuse`, `SubscriptionReuse`, `ResourceGroupReuse`, and `NoReuse`. Defaults to `NoReuse`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Cloud Hardware Security Module Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cloud Hardware Security Module Cluster. The only supported value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cloud Hardware Security Module Cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cloud Hardware Security Module Cluster.

* `activation_state` - The activation state of the Cloud Hardware Security Module Cluster.

* `status_message` - The status message of the Cloud Hardware Security Module Cluster.

* `hsms` - A list of `hsms` blocks as defined below.

* `private_endpoint_connections` - A list of `private_endpoint_connections` blocks as defined below.

---

A `hsms` block exports the following:

* `fqdn` - The fully qualified domain name of the HSM instance.

* `state` - The state of the HSM instance.

* `state_message` - The state message of the HSM instance.

---

A `private_endpoint_connections` block exports the following:

* `id` - The ID of the private endpoint connection.

* `name` - The name of the private endpoint connection.

* `type` - The type of the private endpoint connection.

* `group_ids` - A list of group IDs for the private endpoint connection.

* `private_endpoint` - A `private_endpoint` block as defined below.

* `private_link_service_connection_state` - A `private_link_service_connection_state` block as defined below.

---

A `private_endpoint` block exports the following:

* `id` - The ID of the private endpoint.

---

A `private_link_service_connection_state` block exports the following:

* `status` - The status of the private link service connection.

* `description` - The description of the private link service connection state.

* `actions_required` - The actions required for the private link service connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Cloud Hardware Security Module Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud Hardware Security Module Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the Cloud Hardware Security Module Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the Cloud Hardware Security Module Cluster.

## Import

Cloud Hardware Security Module Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cloud_hardware_security_module_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/cluster1
```
